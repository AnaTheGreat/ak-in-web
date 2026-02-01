package config

import (
	"log"
	"os"
)

type Config struct {
	AppName string
	Env     string
	Port    string
}

func Load() *Config {
	cfg := &Config{
		AppName: getEnv("APP_NAME", "ak-in-web"),
		Env:     getEnv("APP_ENV", "local"),
		Port:    getEnv("APP_PORT", "8080"),
	}

	log.Printf("config loaded: app=%s env=%s", cfg.AppName, cfg.Env)
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
