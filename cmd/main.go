// Package main presents main part that starts application.
package main

import (
	"log"

	"github.com/cyberdr0id/cv-web-service/internal/handler"
	"github.com/cyberdr0id/cv-web-service/internal/server"
)

func main() {
	srv := new(server.Server)

	if err := srv.Run("8000", handler.InitRoutes()); err != nil {
		log.Fatalf("error while starting server: %s", err.Error())
	}
}
