// In server/handlers/collections.go, modify GetCollection to:
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
