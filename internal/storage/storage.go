// Package storage implements functions and types for work with AWS S3 object storage.
package storage

import (
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const urlExpireTime = 5 * time.Minute

// Storage presents type for work with object storage.
type Storage struct {
	storage *s3.S3
	config  *StorageConfig
}

// StorageConfig consist of key parameters for work with object storage.
type StorageConfig struct {
	Bucket      string
	AccessKey   string
	AccessKeyID string
	Region      string
}

// NewStorage creates a new inastance of Storage.
func NewStorage(config *StorageConfig) (*Storage, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.AccessKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot start session: %w", err)
	}

	return &Storage{
		storage: s3.New(session),
		config:  config,
	}, nil
}

// UploadFileToStorage uploads file to object storage.
func (s *Storage) UploadFileToStorage(file io.ReadSeeker, fileID string) error {
	_, err := s.storage.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(fileID),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		return fmt.Errorf("cannot load file to object storage: %w", err)
	}

	return nil
}

// DownloadFileFromStorage downloads file from object storage
func (s *Storage) DownloadFileFromStorage(fileID string) (io.ReadCloser, error) {
	resp, err := s.storage.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot download file: %w", err)
	}

	return resp.Body, nil
}

// GetFileURLByID returns file URL by it ID.
func (s *Storage) GetFileURLByID(fileID string) (string, error) {
	req, _ := s.storage.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(fileID),
	})

	url, err := req.Presign(urlExpireTime)
	if err != nil {
		return "", fmt.Errorf("cannot create file URL: %w", err)
	}

	return url, err
}
