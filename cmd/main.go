// Package main presents main part that starts application.
package main

import (
	"log"

	"github.com/cyberdr0id/referral/internal/handler"
	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/service"
	"github.com/cyberdr0id/referral/internal/storage"
	"github.com/cyberdr0id/referral/pkg/jwt"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error while reading config: %s", err.Error())
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
		log.Fatalf("error while trying to connect to database: %s", err.Error())
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
		log.Fatal(err)
	}

	authService := service.NewAuthService(repo, tm)
	referralService := service.NewReferralService(repo, s3)

	server := handler.NewServer(authService, referralService)

	if err := server.Run(viper.GetString("port"), server); err != nil {
		log.Fatalf("error while starting server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("docs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
