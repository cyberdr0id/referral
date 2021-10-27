// Package cfg contains all configuration data for work with database, object storage, etc.
package cfg

// DatabaseConfig represents a type that contains database configuration data.
type DatabaseConfig struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
	SSLMode      string
}
