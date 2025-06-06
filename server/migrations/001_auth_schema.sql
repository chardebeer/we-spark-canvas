-- Add password_hash column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT;

-- Create user_hearts table for tracking hearts
CREATE TABLE IF NOT EXISTS user_hearts (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    image_id INTEGER NOT NULL REFERENCES images(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, image_id)
);

-- Create index on user_hearts for fast lookup
CREATE INDEX IF NOT EXISTS idx_user_hearts_user_id ON user_hearts(user_id);
CREATE INDEX IF NOT EXISTS idx_user_hearts_image_id ON user_hearts(image_id);

-- Create index on images.tags for fast filtering
CREATE INDEX IF NOT EXISTS idx_images_tags ON images USING GIN(tags);

-- Create index on images.uploaded_at for trending calculations
CREATE INDEX IF NOT EXISTS idx_images_uploaded_at ON images(uploaded_at);