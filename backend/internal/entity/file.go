package entity

import "time"

// File represents a stored object in MinIO/S3.
type File struct {
	ID          string    `json:"id"`
	FileName    string    `json:"file_name"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `json:"mime_type"`
	StoragePath string    `json:"storage_path"`
	UploadedBy  string    `json:"uploaded_by"`
	CreatedAt   time.Time `json:"created_at"`

	URL string `json:"url,omitempty"` // presigned download URL
}

type PresignedURL struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}
