package config

import (
	"log"
	"os"
)

type Config struct {
	// App
	AppName string
	Env     string
	Port    string

	// Database
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func Load() *Config {
	cfg := &Config{
		// App
		AppName: getEnv("APP_NAME", "ak-in-web"),
		Env:     getEnv("APP_ENV", "local"),
		Port:    getEnv("APP_PORT", "8080"),

		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "akdb"),
		DBUser:     getEnv("DB_USER", "akuser"),
		DBPassword: getEnv("DB_PASSWORD", "akpass"),
	}

	log.Printf(
		"config loaded: app=%s env=%s db=%s:%s/%s",
		cfg.AppName, cfg.Env, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
