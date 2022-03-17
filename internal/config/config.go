package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Database presents a database config variables.
type Database struct {
	Host         string `envconfig:"DB_HOST"`
	User         string `envconfig:"DB_USER"`
	Password     string `envconfig:"DB_PASSWORD"`
	DatabaseName string `envconfig:"DB_NAME"`
	Port         string `envconfig:"DB_PORT"`
	SSLMode      string `envconfig:"DB_SSLMODE"`
}

// JWT a type consist of necessary variables for creating JWT token.
type JWT struct {
	Key        string `envconfig:"JWT_KEY"`
	ExpiryTime string `envconfig:"JWT_EXPIRY_TIME"`
}

// GCS consist of necessary parameters for work with Google object storage.
type GCS struct {
	Bucket          string `envconfig:"GOOGLE_BUCKET_NAME"`
	CredentialsPath string `envconfig:"GCS_CREDENTIALS_PATH"`
}

// Application has parameters for starting application.
type Application struct {
	Port string `envconfig:"APP_PORT"`
}

// MainConfig unites all configs of application.
type MainConfig struct {
	DB  *Database
	JWT *JWT
	App *Application
	GCS *GCS
}

// Load loads application config from environment variables.
func Load() (MainConfig, error) {
	var c MainConfig

	err := envconfig.Process("", &c)

	return c, err
}
