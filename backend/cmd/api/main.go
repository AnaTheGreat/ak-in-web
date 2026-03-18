package main

import (
	"log"
	"net/http"

	"ak-in-web/internal"
	"ak-in-web/internal/config"
	"ak-in-web/internal/db"
	"ak-in-web/internal/telemetry"
)

func main() {
	cfg := config.Load()

	shutdownTelemetry := telemetry.Init(cfg.AppName, cfg.Env)
	defer shutdownTelemetry()

	database := db.New(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
	)

	router := internal.Router(database)

	log.Printf("starting server on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))

}
