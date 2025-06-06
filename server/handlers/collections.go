package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pq "github.com/lib/pq"

	"github.com/chardebeer/we-spark-canvas/server/models"
)

// CreateCollection handles POST /collections
func CreateCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
			CreatedBy   int    `json:"created_by" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var id int
		err := db.QueryRow(`INSERT INTO collections (title, description, created_by) VALUES ($1,$2,$3) RETURNING id`,
			input.Title, input.Description, input.CreatedBy).Scan(&id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

// AddImageToCollection handles POST /collections/:id/images
func AddImageToCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		colID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collection id"})
			return
		}
		var input struct {
			ImageID int `json:"image_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_, err = db.Exec(`INSERT INTO collection_images (collection_id, image_id) VALUES ($1,$2)`, colID, input.ImageID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	}
}

// GetCollections handles GET /collections
func GetCollections(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, title, description, created_by FROM collections ORDER BY id DESC`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var cols []models.Collection
		for rows.Next() {
			var col models.Collection
			if err := rows.Scan(&col.ID, &col.Title, &col.Description, &col.CreatedBy); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			cols = append(cols, col)
		}
		c.JSON(http.StatusOK, cols)
	}
}

// GetCollection handles GET /collections/:id and returns collection with its images
func GetCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		colIDStr := c.Param("id")
		colID, err := strconv.Atoi(colIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collection id"})
			return
		}
		var col models.Collection
		row := db.QueryRow(`SELECT id, title, description, created_by FROM collections WHERE id=$1`, colID)
		if err := row.Scan(&col.ID, &col.Title, &col.Description, &col.CreatedBy); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		rows, err := db.Query(
			`SELECT images.id, images.url, images.caption, images.tags, images.hearts,
                    images.uploaded_by, images.uploaded_at
             FROM images
             JOIN collection_images ON images.id = collection_images.image_id
             WHERE collection_images.collection_id = $1`,
			colID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var imgs []models.Image
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
			imgs = append(imgs, img)
		}
		c.JSON(http.StatusOK, gin.H{
			"collection": col,
			"images":     imgs,
		})
	}
}
