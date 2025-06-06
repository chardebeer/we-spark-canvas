package handlers

import (
  "database/sql"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
  "github.com/chardebeer/we-spark-canvas/server/models"
)

// CreateUser handles POST /users
func CreateUser(db *sql.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    var input struct {
      Username  string `json:"username" binding:"required"`
      AvatarURL string `json:"avatar_url"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
    }
    var userID int
    err := db.QueryRow(
      "INSERT INTO users (username, avatar_url) VALUES ($1, $2) RETURNING id",
      input.Username, input.AvatarURL,
    ).Scan(&userID)
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    c.JSON(http.StatusCreated, gin.H{"id": userID})
  }
}

// GetUser handles GET /users/:id
func GetUser(db *sql.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
      return
    }
    var user models.User
    row := db.QueryRow("SELECT id, username, avatar_url FROM users WHERE id=$1", id)
    if err := row.Scan(&user.ID, &user.Username, &user.AvatarURL); err != nil {
      if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
      } else {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      }
      return
    }
    c.JSON(http.StatusOK, user)
  }
}
// GetUserImages handles GET /users/:id/images
func GetUserImages(db *sql.DB) gin.HandlerFunc {
  return func(c *gin.Context) {
    idStr := c.Param("id")
    userID, err := strconv.Atoi(idStr)
    if err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
      return
    }
    rows, err := db.Query(
      `SELECT id, url, caption, tags, hearts, uploaded_by, uploaded_at
       FROM images
       WHERE uploaded_by = $1
       ORDER BY uploaded_at DESC`,
      userID,
    )
    if err != nil {
      c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
      return
    }
    defer rows.Close()

    var images []models.Image
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
