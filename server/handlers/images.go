package handlers

import (
  "database/sql"                // for *sql.DB, sql.ErrNoRows
  "net/http"
  "strconv"
  "strings"

  "github.com/gin-gonic/gin"
  shell "github.com/ipfs/go-ipfs-api"
  pq "github.com/lib/pq"         // ← Import pq for pq.Array usage
  "github.com/chardebeer/we-spark-canvas/server/models"
)

// UploadImage handles POST /images
// Expects multipart form: file, caption, tags (comma-separated)
func UploadImage(db *sql.DB, ipfsShell *shell.Shell) gin.HandlerFunc {
  return func(c *gin.Context) {
    // Get user ID from context (set by AuthMiddleware)
    userID, exists := c.Get("userID")
    if !exists {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
      return
    }

    // 1. Parse multipart form
    file, err := c.FormFile("file")
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
      return
    }
    caption := c.PostForm("caption")
    tagsStr := c.PostForm("tags")

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
      url, caption, pq.Array(tagsArray), userID,
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
    // Get user ID from context (set by AuthMiddleware)
    userID, exists := c.Get("userID")
    if !exists {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
      return
    }

    idStr := c.Param("id")
    imageID, err := strconv.Atoi(idStr)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "invalid image id"})
      return
    }

    // Check if the user has already hearted this image
    var alreadyExists bool
    err = db.QueryRow(
      `SELECT EXISTS(SELECT 1 FROM user_hearts WHERE user_id = $1 AND image_id = $2)`,
      userID, imageID,
    ).Scan(&alreadyExists)
    
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    if alreadyExists {
      c.JSON(http.StatusConflict, gin.H{"error": "user has already hearted this image"})
      return
    }

    // Begin transaction
    tx, err := db.Begin()
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    defer tx.Rollback()

    // Insert into user_hearts
    _, err = tx.Exec(
      `INSERT INTO user_hearts (user_id, image_id, created_at) VALUES ($1, $2, NOW())`,
      userID, imageID,
    )
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    // Increment heart count
    _, err = tx.Exec(`UPDATE images SET hearts = hearts + 1 WHERE id = $1`, imageID)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    // Get the new heart count
    var newHeartCount int
    err = tx.QueryRow(`SELECT hearts FROM images WHERE id = $1`, imageID).Scan(&newHeartCount)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }

    c.JSON(http.StatusOK, gin.H{
      "message": "hearted successfully",
      "hearts": newHeartCount,
    })
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

// GetImagesByTag handles GET /images/tag/:tag?limit=20&offset=0
func GetImagesByTag(db *sql.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    tag := c.Param("tag")
    if tag == "" {
      c.JSON(http.StatusBadRequest, gin.H{"error": "tag parameter is required"})
      return
    }

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
       FROM images 
       WHERE $1 = ANY(tags)
       ORDER BY uploaded_at DESC LIMIT $2 OFFSET $3`,
      tag, limit, offset,
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

// TagCount represents a tag and its usage count
type TagCount struct {
  Tag   string `json:"tag"`
  Count int    `json:"count"`
}

// GetTrendingTags handles GET /tags/trending
func GetTrendingTags(db *sql.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    // Get top 20 tags by usage count in the last 24 hours
    rows, err := db.Query(`
      SELECT unnest(tags) as tag, COUNT(*) as count
      FROM images
      WHERE uploaded_at > NOW() - INTERVAL '24 hours'
      GROUP BY tag
      ORDER BY count DESC
      LIMIT 20
    `)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    defer rows.Close()

    tags := []TagCount{}
    for rows.Next() {
      var tag TagCount
      if err := rows.Scan(&tag.Tag, &tag.Count); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
      }
      tags = append(tags, tag)
    }

    // If no trending tags in the last 24 hours, get all-time trending tags
    if len(tags) == 0 {
      rows, err = db.Query(`
        SELECT unnest(tags) as tag, COUNT(*) as count
        FROM images
        GROUP BY tag
        ORDER BY count DESC
        LIMIT 20
      `)
      if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
      }
      defer rows.Close()

      for rows.Next() {
        var tag TagCount
        if err := rows.Scan(&tag.Tag, &tag.Count); err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
          return
        }
        tags = append(tags, tag)
      }
    }

    c.JSON(http.StatusOK, tags)
  }
}
