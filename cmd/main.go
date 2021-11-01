// Package main presents main part that starts application.
package main

import (
	"log"

	"github.com/cyberdr0id/referral/internal/handler"
	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/server"
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
	_ = repo

	srv := server.NewServer(viper.GetString("port"), handler.InitRoutes())

	if err := srv.Run(); err != nil {
		log.Fatalf("error while starting server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("docs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
