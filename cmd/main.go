package main

import (
	"log"

	"github.com/cyberdr0id/cv-web-service/internal/handler"
	"github.com/cyberdr0id/cv-web-service/internal/server"
)

func main() {
	srv := new(server.Server)
	handlers := new(handler.Handler)

	if err := srv.Run("8000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error while starting server: %s", err.Error())
	}
}
