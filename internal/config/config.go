package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Database struct {
	Host         string `envconfig:"DB_HOST"`
	User         string `envconfig:"DB_USER"`
	Password     string `envconfig:"DB_PASSWORD"`
	DatabaseName string `envconfig:"DB_NAME"`
	Port         string `envconfig:"DB_PORT"`
	SSLMode      string `envconfig:"DB_SSLMODE"`
}

type JWT struct {
	Key        string `envconfig:"JWT_KEY"`
	ExpiryTime string `envconfig:"JWT_EXPIRY_TIME"`
}

type AWS struct {
	Bucket      string `envconfig:"AWS_BUCKET"`
	Region      string `envconfig:"AWS_REGION"`
	AccessKey   string `envconfig:"AWS_ACCESS_KEY"`
	AccessKeyID string `envconfig:"AWS_ACCESS_KEY_ID"`
}

type Application struct {
	Port string `envconfig:"APP_PORT"`
}

type MainConfig struct {
	DB  *Database
	JWT *JWT
	AWS *AWS
	App *Application
}

func Load() (MainConfig, error) {
	var c MainConfig

	err := envconfig.Process("", &c)

	return c, err
}
