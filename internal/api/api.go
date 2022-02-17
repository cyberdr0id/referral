package api

import (
	"fmt"
	"log"
	"os"

	"github.com/cyberdr0id/referral/internal/handler"
	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/service"
	"github.com/cyberdr0id/referral/internal/storage"
	"github.com/cyberdr0id/referral/pkg/jwt"
	mylog "github.com/cyberdr0id/referral/pkg/log"
)

// Start starts API with initialization of necessary components.
func Start() (*mylog.Logger, error) {
	logger, err := mylog.NewLogger()
	if err != nil {
		log.Fatalf("error with logger creating: %s", err.Error())
	}

	config := repository.DatabaseConfig{
		Host:         os.Getenv("DB_HOST"),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		Port:         os.Getenv("DB_PORT"),
		SSLMode:      os.Getenv("DB_SSLMODE"),
	}

	db, err := repository.NewConnection(config)
	if err != nil {
		return logger, fmt.Errorf("error while trying to connect to database: %s", err)
	}

	repo := repository.NewRepository(db)

	tm, err := jwt.NewTokenManager(
		os.Getenv("JWT_KEY"),
		os.Getenv("JWT_EXPIRY_TIME"),
	)
	if err != nil {
		return logger, fmt.Errorf("error with creating JWT token manager: %w", err)
	}

	s3config := &storage.StorageConfig{
		Bucket:      os.Getenv("AWS_BUCKET"),
		Region:      os.Getenv("AWS_REGION"),
		AccessKey:   os.Getenv("AWS_ACCESS_KEY"),
		AccessKeyID: os.Getenv("AWS_ACCESS_KEY_ID"),
	}

	s3, err := storage.NewStorage(s3config)
	if err != nil {
		return logger, fmt.Errorf("cannot create new instance of object storage: %s", err)
	}

	authService := service.NewAuthService(repo, tm)
	referralService := service.NewReferralService(repo, s3)

	server := handler.NewServer(authService, referralService, logger)

	if err := server.Run(os.Getenv("APP_PORT"), server); err != nil {
		return logger, fmt.Errorf("error while starting server: %s", err)
	}

	return logger, nil
}
