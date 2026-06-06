package entity

import "time"

type Attachment struct {
	ID          string     `json:"id"`
	ParentType  string     `json:"parent_type"` // issue | page | comment
	ParentID    string     `json:"parent_id"`
	FileName    string     `json:"file_name"`
	FileSize    int64      `json:"file_size"`
	MimeType    string     `json:"mime_type"`
	StoragePath string     `json:"storage_path"`
	StorageType string     `json:"storage_type"` // local | s3
	Checksum    *string    `json:"checksum,omitempty"`
	UploadedBy  string     `json:"uploaded_by"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"-"`

	Uploader   *UserShort `json:"uploader,omitempty"`
	DownloadURL string    `json:"download_url,omitempty"`
}
