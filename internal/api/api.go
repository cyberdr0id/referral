package api

import (
	"fmt"
	"log"

	"github.com/cyberdr0id/referral/internal/handler"
	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/service"
	"github.com/cyberdr0id/referral/internal/storage"
	"github.com/cyberdr0id/referral/pkg/jwt"
	mylog "github.com/cyberdr0id/referral/pkg/log"
	"github.com/spf13/viper"
)

func initConfig() error {
	viper.AddConfigPath("docs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

// Start starts API with initialization of necessary components.
func Start() (*mylog.Logger, error) {
	logger, err := mylog.NewLogger()
	if err != nil {
		log.Fatalf("error with logger creating: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		return &mylog.Logger{}, fmt.Errorf("error while reading config: %s", err.Error())
	}

	config := repository.DatabaseConfig{
		Host:         viper.GetString("db.host"),
		User:         viper.GetString("db.user"),
		Password:     viper.GetString("db.password"),
		DatabaseName: viper.GetString("db.dbname"),
		Port:         viper.GetString("db.port"),
		SSLMode:      viper.GetString("db.sslmode"),
	}

	db, err := repository.NewConnection(config)
	if err != nil {
		return &mylog.Logger{}, fmt.Errorf("error while trying to connect to database: %s", err)
	}

	repo := repository.NewRepository(db)
	tm := jwt.NewTokenManager(viper.GetString("jwt.key"), viper.GetInt("jwt.expiryTime"))

	s3config := &storage.StorageConfig{
		Bucket:      viper.GetString("aws.bucket"),
		Region:      viper.GetString("aws.region"),
		AccessKey:   viper.GetString("aws.accessKey"),
		AccessKeyID: viper.GetString("aws.accessKeyID"),
	}

	s3, err := storage.NewStorage(s3config)
	if err != nil {
		return &mylog.Logger{}, fmt.Errorf("cannot create new instance of object storage: %s", err)
	}

	authService := service.NewAuthService(repo, tm)
	referralService := service.NewReferralService(repo, s3)

	server := handler.NewServer(authService, referralService, logger)

	if err := server.Run(viper.GetString("port"), server); err != nil {
		return &mylog.Logger{}, fmt.Errorf("error while starting server: %s", err)
	}

	return logger, nil
}
