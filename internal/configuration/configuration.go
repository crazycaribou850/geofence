package configuration

import (
	"os"
)

const AppName = "geofence"

type Config struct {
	DBURL string
	Port string
}

func Load() *Config {
	dbURL := loadPSQLConfig()
	port := loadHTTPConfig()
	return &Config{dbURL, port}
}

func loadPSQLConfig() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:@localhost:5432/geofence?ssl=false"
	}
	return dbURL
}

func loadHTTPConfig() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}
	return port
}
