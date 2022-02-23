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

// AWS is a struct of AWS S3 configuration.
type AWS struct {
	Bucket      string `envconfig:"AWS_BUCKET"`
	Region      string `envconfig:"AWS_REGION"`
	AccessKey   string `envconfig:"AWS_ACCESS_KEY"`
	AccessKeyID string `envconfig:"AWS_ACCESS_KEY_ID"`
}

// Application has parameters for starting application.
type Application struct {
	Port string `envconfig:"APP_PORT"`
}

// MainConfig unites all configs of application.
type MainConfig struct {
	DB  *Database
	JWT *JWT
	AWS *AWS
	App *Application
}

// Load loads application config from environment variables.
func Load() (MainConfig, error) {
	var c MainConfig

	err := envconfig.Process("", &c)

	return c, err
}
