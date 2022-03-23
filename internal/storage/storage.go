// Package storage implements functions and types for work with AWS S3 object storage.
package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/cyberdr0id/referral/internal/config"
	"google.golang.org/api/option"
)

const (
	pdfType              = "application/pdf"
	maxFileSize          = 32 << 20
	optionExpirationTime = 15
	getMethod            = "GET"
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
func (s *Storage) DownloadFile(ctx context.Context, fileID, fileName string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  getMethod,
		Expires: time.Now().Add(optionExpirationTime * time.Minute),
	}

	u, err := s.Client.Bucket(s.Bucket).SignedURL(fileName, opts)
	if err != nil {
		return "", fmt.Errorf("cannot created file URL: %w", err)
	}

	return u, nil
}

// UploadFileToStorage uploads file to object storage.
func (s *Storage) UploadFileToStorage(ctx context.Context, file multipart.File, fileID string) error {
	wc := s.Client.Bucket(s.Bucket).Object(fileID).NewWriter(ctx)

	wc.Size = maxFileSize
	wc.ContentType = pdfType

	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("cannot copy file info: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("cannot close file writer: %w", err)
	}

	return nil
}
