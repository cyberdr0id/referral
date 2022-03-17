// Package storage implements functions and types for work with AWS S3 object storage.
package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/cyberdr0id/referral/internal/config"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"os"
)

// Storage presents a type for work with Google cloud storage
type Storage struct {
	Client *storage.Client
	Bucket string
}

// NewStorage creates a new instance of Storage.
func NewStorage(cfg *config.GCS) (*Storage, error) {
	newClient, err := storage.NewClient(context.Background(), option.WithCredentialsFile(cfg.CredentialsPath))
	if err != nil {
		return &Storage{}, fmt.Errorf("cannot create new client of object storage: %w", err)
	}

	return &Storage{
		Client: newClient,
		Bucket: cfg.Bucket,
	}, nil
}

// DownloadFile downloads a file from object storage by file id.
func (s *Storage) DownloadFile(ctx context.Context, fileID, fileName string) error {
	f, err := os.Create(fmt.Sprintf("downloaded/%s.pdf", fileName))
	if err != nil {
		return fmt.Errorf("cannot create file: %w", err)
	}

	rc, err := s.Client.Bucket(s.Bucket).Object(fileID).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("cannot read object: %w", err)
	}

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("cannot copy object to file: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("cannot close file: %w", err)
	}
	if err = rc.Close(); err != nil {
		return fmt.Errorf("cannot close file reader: %w", err)
	}

	return nil
}

// UploadFileToStorage uploads file to object storage.
func (s *Storage) UploadFileToStorage(ctx context.Context, file multipart.File, fileID string) error {
	wc := s.Client.Bucket(s.Bucket).Object(fileID).NewWriter(ctx)

	wc.Size = 32 << 20
	wc.ContentType = "application/pdf"

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("cannot copy file info: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("cannot close file writer: %w", err)
	}

	return nil
}
