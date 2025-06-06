package storage

import (
  "database/sql"
  _ "github.com/lib/pq"
  "log"
  "os"
)

// NewPostgresDB reads DATABASE_URL from the environment and returns a live DB connection.
func NewPostgresDB() *sql.DB {
  dsn := os.Getenv("DATABASE_URL")
  db, err := sql.Open("postgres", dsn)
  if err != nil {
    log.Fatalf("failed to connect to Postgres: %v", err)
  }
  if err := db.Ping(); err != nil {
    log.Fatalf("failed to ping DB: %v", err)
  }
  return db
}
