// Package storage implements functions and types for work with AWS S3 object storage.
package storage

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

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
	session, err := startSession(config)
	if err != nil {
		return nil, fmt.Errorf("cannot start session: %w", err)
	}

	return &Storage{
		storage: s3.New(session),
		config:  config,
	}, nil
}

// startSession creates a new storage session.
func startSession(config *StorageConfig) (*session.Session, error) {
	return session.NewSession(&aws.Config{
		Region:      aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.AccessKey, ""),
	})
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
func (s *Storage) DownloadFileFromStorage(fileID string) (io.ReadSeeker, error) {
	resp, err := s.storage.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, fmt.Errorf("cannot download file: %w", err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read downloaded data: %w", err)
	}

	return bytes.NewReader(buf), nil
}
