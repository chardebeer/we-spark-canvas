package handlers

import (
	"database/sql" // for *sql.DB, sql.ErrNoRows
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
	pq "github.com/lib/pq" // ← Import pq for pq.Array usage

	"github.com/chardebeer/we-spark-canvas/server/models"
)

// UploadImage handles POST /images
// Expects multipart form: file, caption, tags (comma-separated), uploaded_by
func UploadImage(db *sql.DB, ipfsShell *shell.Shell) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Parse multipart form
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}
		caption := c.PostForm("caption")
		tagsStr := c.PostForm("tags")
		uploadedByStr := c.PostForm("uploaded_by")
		uploadedBy, err := strconv.Atoi(uploadedByStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uploaded_by"})
			return
		}

		// 2. Read file into IPFS
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()

		cid, err := ipfsShell.Add(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "IPFS upload failed"})
			return
		}
		// Construct a gateway URL for retrieval
		url := "https://ipfs.io/ipfs/" + cid

		// 3. Insert into PostgreSQL
		tagsArray := []string{}
		if tagsStr != "" {
			tagsArray = strings.Split(tagsStr, ",")
		}
		var imageID int
		err = db.QueryRow(
			`INSERT INTO images (url, caption, tags, hearts, uploaded_by) 
       VALUES ($1, $2, $3, 0, $4) RETURNING id`,
			url, caption, pq.Array(tagsArray), uploadedBy,
		).Scan(&imageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": imageID, "url": url})
	}
}

// GetImages handles GET /images?limit=20&offset=0
func GetImages(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 20
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}

		rows, err := db.Query(
			`SELECT id, url, caption, tags, hearts, uploaded_by, uploaded_at 
       FROM images ORDER BY uploaded_at DESC LIMIT $1 OFFSET $2`,
			limit, offset,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		images := []models.Image{}
		for rows.Next() {
			var img models.Image
			var tags pq.StringArray
			if err := rows.Scan(
				&img.ID, &img.URL, &img.Caption, &tags, &img.Hearts, &img.UploadedBy, &img.UploadedAt,
			); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			img.Tags = tags
			images = append(images, img)
		}
		c.JSON(http.StatusOK, images)
	}
}

// HeartImage handles POST /images/:id/heart
func HeartImage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		imageID, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image id"})
			return
		}
		var input struct {
			UserID int `json:"user_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(`INSERT INTO image_hearts (image_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, imageID, input.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		_, err = tx.Exec(`UPDATE images SET hearts = (SELECT COUNT(*) FROM image_hearts WHERE image_id=$1) WHERE id=$1`, imageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var count int
		if err := tx.QueryRow(`SELECT hearts FROM images WHERE id=$1`, imageID).Scan(&count); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"hearts": count})
	}
}

// GetImage handles GET /images/:id
func GetImage(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image id"})
			return
		}
		var img models.Image
		var tags pq.StringArray
		row := db.QueryRow(
			`SELECT id, url, caption, tags, hearts, uploaded_by, uploaded_at 
       FROM images WHERE id=$1`, id,
		)
		if err := row.Scan(
			&img.ID, &img.URL, &img.Caption, &tags, &img.Hearts, &img.UploadedBy, &img.UploadedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		img.Tags = tags
		c.JSON(http.StatusOK, img)
	}
}

// GetImagesByTag handles GET /images/tag/:tag
func GetImagesByTag(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := c.Param("tag")
		limitStr := c.DefaultQuery("limit", "20")
		offsetStr := c.DefaultQuery("offset", "0")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			limit = 20
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			offset = 0
		}
		rows, err := db.Query(`SELECT id, url, caption, tags, hearts, uploaded_by, uploaded_at FROM images WHERE $1 = ANY(tags) ORDER BY uploaded_at DESC LIMIT $2 OFFSET $3`, tag, limit, offset)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var images []models.Image
		for rows.Next() {
			var img models.Image
			var tags pq.StringArray
			if err := rows.Scan(&img.ID, &img.URL, &img.Caption, &tags, &img.Hearts, &img.UploadedBy, &img.UploadedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			img.Tags = tags
			images = append(images, img)
		}
		c.JSON(http.StatusOK, images)
	}
}

// GetTrendingTags handles GET /tags/trending
func GetTrendingTags(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`SELECT tag, COUNT(*) as cnt FROM (SELECT unnest(tags) AS tag FROM images WHERE uploaded_at >= NOW() - INTERVAL '24 hours') t GROUP BY tag ORDER BY cnt DESC LIMIT 20`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		type TagCount struct {
			Tag   string `json:"tag"`
			Count int    `json:"count"`
		}
		var tags []TagCount
		for rows.Next() {
			var tc TagCount
			if err := rows.Scan(&tc.Tag, &tc.Count); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			tags = append(tags, tc)
		}
		c.JSON(http.StatusOK, tags)
	}
}
