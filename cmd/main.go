// Package main presents main part that starts application.
package main

import (
	"github.com/cyberdr0id/referral/internal/api"
)

func main() {
	if logger, err := api.Start(); err != nil {
		logger.ErrorLogger.Fatal(err)
	}
}
