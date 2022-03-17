package api

import (
	"fmt"
	"github.com/cyberdr0id/referral/internal/config"
	"github.com/cyberdr0id/referral/internal/handler"
	"github.com/cyberdr0id/referral/internal/repository"
	"github.com/cyberdr0id/referral/internal/service"
	"github.com/cyberdr0id/referral/internal/storage"
	"github.com/cyberdr0id/referral/pkg/jwt"
	mylog "github.com/cyberdr0id/referral/pkg/log"
	"log"
)

// Start starts API with initialization of necessary components.
func Start() (*mylog.Logger, error) {
	logger, err := mylog.NewLogger()
	if err != nil {
		log.Fatalf("error with logger creating: %s", err.Error())
	}

	cfg, err := config.Load()
	if err != nil {
		return logger, fmt.Errorf("cannot read application config: %w", err)
	}

	db, err := repository.NewConnection(cfg.DB)
	if err != nil {
		return logger, fmt.Errorf("error while trying to connect to database: %s", err)
	}

	repo := repository.NewRepository(db)

	tm, err := jwt.NewTokenManager(cfg.JWT)
	if err != nil {
		return logger, fmt.Errorf("error with creating JWT token manager: %w", err)
	}

	gcs, err := storage.NewStorage(cfg.GCS)
	if err != nil {
		return logger, fmt.Errorf("cannot create new instance of object storage: %s", err)
	}

	authService := service.NewAuthService(repo, tm)
	referralService := service.NewReferralService(repo, gcs)

	server := handler.NewServer(authService, referralService, logger)

	if err := server.Run(cfg.App.Port, server); err != nil {
		return logger, fmt.Errorf("error while starting server: %s", err)
	}

	return logger, nil
}
