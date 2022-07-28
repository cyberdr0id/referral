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
	"github.com/kelseyhightower/envconfig"
)

type appConfig struct {
	Port string `envconfig:"APP_PORT"`
}

// Start starts API with initialization of necessary components.
func Start() (*mylog.Logger, error) {
	logger, err := mylog.NewLogger()
	if err != nil {
		log.Fatalf("error with logger creating: %s", err.Error())
	}

	db, err := repository.NewConnection()
	if err != nil {
		fmt.Println(fmt.Errorf("error while trying to connect to database: %s", err))
		return logger, fmt.Errorf("error while trying to connect to database: %s", err)
	}

	repo := repository.NewRepository(db)

	tm, err := jwt.NewTokenManager()
	if err != nil {
		fmt.Println(fmt.Errorf("error with creating JWT token manager: %w", err))
		return logger, fmt.Errorf("error with creating JWT token manager: %w", err)
	}

	gcs, err := storage.NewStorage()
	if err != nil {
		fmt.Println(fmt.Errorf("cannot create new instance of object storage: %s", err))
		return logger, fmt.Errorf("cannot create new instance of object storage: %s", err)
	}

	authService := service.NewAuthService(repo, tm)
	referralService := service.NewReferralService(repo, gcs)

	server := handler.NewServer(authService, referralService, logger)

	cfg, err := loadConfig()
	if err != nil {
		return logger, fmt.Errorf("error with loading app config: %w", err)
	}
	log.Println(cfg)
	if err := server.Run(cfg.Port, server); err != nil {
		fmt.Println(fmt.Errorf("error while starting server: %s", err))
		return logger, fmt.Errorf("error while starting server: %s", err)
	}

	return logger, nil
}

func loadConfig() (*appConfig, error) {
	var c appConfig

	if err := envconfig.Process("app", &c); err != nil {
		return nil, fmt.Errorf("unable to read application config: %w", err)
	}

	return &c, nil
}
