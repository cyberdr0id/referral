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
	options := fmt.Sprintf(`{
		"type": "%s",
		"project_id": "%s",
		"private_key_id": "%s",
		"private_key": "-----BEGIN PRIVATE KEY-----\n%s\n-----END PRIVATE KEY-----\n",
		"client_email": "%s",
		"client_id": "%s",
		"auth_uri": "%s",
		"token_uri": "%s",
		"auth_provider_x509_cert_url": "%s",
		"client_x509_cert_url": "%s"
	}`, cfg.Type, cfg.ProjectID, cfg.PrivateKeyID, cfg.PrivateKey, cfg.ClientEmail, cfg.ClientID, cfg.AuthURI, cfg.TokenURI, cfg.AuthProvider, cfg.ClientURL)

	newClient, err := storage.NewClient(context.Background(), option.WithCredentialsJSON([]byte(options)))
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
