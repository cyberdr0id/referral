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
	"github.com/kelseyhightower/envconfig"
)

const (
	urlExpireTime = 10 * time.Minute
)

type storageConfig struct {
	Bucket      string `envconfig:"AWS_BUCKET"`
	Region      string `envconfig:"AWS_REGION"`
	AccessKey   string `envconfig:"AWS_ACCESS_KEY"`
	AccessKeyID string `envconfig:"AWS_ACCESS_KEY_ID"`
}

type Storage struct {
	cfg *storageConfig
	s3  *s3.S3
}

func NewStorage() (*Storage, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("error with config loading: %w", err)
	}

	sn, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.AccessKey, ""),
	})
	if err != nil {
		return &Storage{}, fmt.Errorf("unable to create session: %w", err)
	}

	return &Storage{
		s3:  s3.New(sn),
		cfg: cfg,
	}, nil
}

func loadConfig() (*storageConfig, error) {
	var c storageConfig

	err := envconfig.Process("aws", &c)
	if err != nil {
		return nil, fmt.Errorf("error with config loading: %w", err)
	}

	return &c, nil
}

func (s *Storage) UploadFile(file io.ReadSeeker, fileName string) error {
	_, err := s.s3.PutObject(&s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(fileName),
		ACL:    aws.String(s3.BucketCannedACLPublicRead),
	})
	if err != nil {
		return fmt.Errorf("unable put file to object storage: %w", err)
	}

	return nil
}

func (s *Storage) DownloadFile(fileID string) (io.ReadCloser, error) {
	resp, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get object with key %s from storage: %w", fileID, err)
	}

	return resp.Body, nil
}

func (s *Storage) GetFileURL(fileID string) (string, error) {
	req, _ := s.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(fileID),
	})

	url, err := req.Presign(urlExpireTime)
	if err != nil {
		return "", fmt.Errorf("unable create request's signed URL: %w", err)
	}

	return url, nil
}
