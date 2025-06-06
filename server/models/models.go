package models

import "time"

// User represents a user in the system.
type User struct {
  ID        int       `json:"id"`
  Username  string    `json:"username"`
  AvatarURL string    `json:"avatar_url,omitempty"`
}

// Image represents an uploaded image.
type Image struct {
  ID         int       `json:"id"`
  URL        string    `json:"url"`
  Caption    string    `json:"caption,omitempty"`
  Tags       []string  `json:"tags,omitempty"`
  Hearts     int       `json:"hearts"`
  UploadedBy int       `json:"uploaded_by"`
  UploadedAt time.Time `json:"uploaded_at"`
}

// Collection represents a user-created collection (moodboard).
type Collection struct {
  ID          int    `json:"id"`
  Title       string `json:"title"`
  Description string `json:"description,omitempty"`
  CreatedBy   int    `json:"created_by"`
}
