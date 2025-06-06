package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	
	"github.com/chardebeer/we-spark-canvas/server/models"
)

// CreateCollectionRequest represents the request body for creating a collection
type CreateCollectionRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// GetAllCollections returns a list of all collections
func GetAllCollections(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rows, err := db.Query(`SELECT id, title, description, created_by FROM collections ORDER BY id DESC`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var collections []models.Collection
		for rows.Next() {
			var col models.Collection
			if err := rows.Scan(&col.ID, &col.Title, &col.Description, &col.CreatedBy); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			collections = append(collections, col)
		}

		if err := rows.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, collections)
	}
}

// CreateCollection creates a new collection
func CreateCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
			return
		}

		var req CreateCollectionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var collectionID int
		err := db.QueryRow(
			`INSERT INTO collections (title, description, created_by) 
			VALUES ($1, $2, $3) RETURNING id`,
			req.Title, req.Description, userID,
		).Scan(&collectionID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"id":          collectionID,
			"title":       req.Title,
			"description": req.Description,
			"created_by":  userID,
		})
	}
}

// AddImageToCollectionRequest represents the request body for adding an image to a collection
type AddImageToCollectionRequest struct {
	ImageID int `json:"image_id" binding:"required"`
}

// AddImageToCollection adds an image to a collection
func AddImageToCollection(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		colIDStr := c.Param("id")
		colID, err := strconv.Atoi(colIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collection id"})
			return
		}

		var req AddImageToCollectionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if collection exists
		var exists bool
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM collections WHERE id = $1)`, colID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "collection not found"})
			return
		}

		// Check if image exists
		err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM images WHERE id = $1)`, req.ImageID).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "image not found"})
			return
		}

		// Check if image is already in collection
		err = db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM collection_images WHERE collection_id = $1 AND image_id = $2)`,
			colID, req.ImageID,
		).Scan(&exists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if exists {
			c.JSON(http.StatusConflict, gin.H{"error": "image already in collection"})
			return
		}

		_, err = db.Exec(
			`INSERT INTO collection_images (collection_id, image_id) VALUES ($1, $2)`,
			colID, req.ImageID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "image added to collection"})
	}
}

// GetCollection returns a single collection and its images
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
