// Package main presents main part that starts application.
package main

import (
	"log"

	"github.com/cyberdr0id/referral/internal/api"
	mylog "github.com/cyberdr0id/referral/pkg/log"
)

func main() {
	logger, err := mylog.NewLogger()
	if err != nil {
		log.Fatalf("error with logger creating: %s", err.Error())
	}

	if err := api.Start(logger); err != nil {
		logger.ErrorLogger.Fatal(err)
	}
}
