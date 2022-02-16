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
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("./")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	return viper.ReadInConfig()
}

// Start starts API with initialization of necessary components.
func Start() (*mylog.Logger, error) {
	logger, err := mylog.NewLogger()
	if err != nil {
		log.Fatalf("error with logger creating: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		return logger, fmt.Errorf("error while reading config: %s", err.Error())
	}

	config := repository.DatabaseConfig{
		Host:         viper.GetString("DB_HOST"),
		User:         viper.GetString("DB_USER"),
		Password:     viper.GetString("DB_PASSWORD"),
		DatabaseName: viper.GetString("DB_NAME"),
		Port:         viper.GetString("DB_PORT"),
		SSLMode:      viper.GetString("DB_SSLMODE"),
	}

	db, err := repository.NewConnection(config)
	if err != nil {
		return logger, fmt.Errorf("error while trying to connect to database: %s", err)
	}

	repo := repository.NewRepository(db)
	tm := jwt.NewTokenManager(
		viper.GetString("JWT_KEY"),
		viper.GetInt("JWT_EXPIRY_TIME"),
	)

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

	if err := server.Run(viper.GetString("APP_PORT"), server); err != nil {
		return logger, fmt.Errorf("error while starting server: %s", err)
	}

	return logger, nil
}
