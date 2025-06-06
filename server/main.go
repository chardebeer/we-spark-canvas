package main

import (
  "log"
  "os"
  "time"

  "github.com/gin-contrib/cors"
  "github.com/gin-gonic/gin"
  "github.com/joho/godotenv"
  shell "github.com/ipfs/go-ipfs-api"

  "github.com/chardebeer/we-spark-canvas/server/handlers"
  "github.com/chardebeer/we-spark-canvas/server/storage"
)

func main() {
  godotenv.Load(".env")

  db := storage.NewPostgresDB()
  ipfsURL := os.Getenv("IPFS_API_URL")
  if ipfsURL == "" {
    ipfsURL = "http://127.0.0.1:5001"
  }
  ipfsShell := shell.NewShell(ipfsURL)

  router := gin.Default()

  // Enable CORS for your Next.js frontend
  router.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
  }))
  
  // Auth routes - public
  router.POST("/auth/register", handlers.Register(db))
  router.POST("/auth/login", handlers.Login(db))

  // Public routes
  router.GET("/users/:id", handlers.GetUser(db))
  router.GET("/images", handlers.GetImages(db))
  router.GET("/images/:id", handlers.GetImage(db))
  router.GET("/images/tag/:tag", handlers.GetImagesByTag(db))
  router.GET("/tags/trending", handlers.GetTrendingTags(db))
  router.GET("/collections/:id", handlers.GetCollection(db))
  router.GET("/collections", handlers.GetAllCollections(db))
  
  // Protected routes - require authentication
  authorized := router.Group("/")
  authorized.Use(handlers.AuthMiddleware())
  {
    // User routes
    authorized.GET("/users/:id/images", handlers.GetUserImages(db))
    authorized.GET("/users/:id/hearts", handlers.GetUserHearts(db))
    
    // Image routes
    authorized.POST("/images", handlers.UploadImage(db, ipfsShell))
    authorized.POST("/images/:id/heart", handlers.HeartImage(db))
    
    // Collection routes
    authorized.POST("/collections", handlers.CreateCollection(db))
    authorized.POST("/collections/:id/images", handlers.AddImageToCollection(db))
  }

  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  log.Printf("🚀 We Spark Canvas REST API running on port %s", port)
  if err := router.Run(":" + port); err != nil {
    log.Fatalf("failed to run server: %v", err)
  }
}
