package main

import (
	"log"
	"net/http"
	"time"

	"ak-in-web/internal"
	"ak-in-web/internal/config"
	"ak-in-web/internal/db"
	"ak-in-web/internal/telemetry"
)

func main() {
	cfg := config.Load()

	shutdownTelemetry := telemetry.Init(cfg.AppName, cfg.Env)
	defer shutdownTelemetry()

	database, err := db.New(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBSSLMode,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err := db.WaitForConnection(database, 30, time.Second); err != nil {
		log.Fatal(err)
	}
	if err := db.RunMigrations(database, "migrations"); err != nil {
		log.Fatal(err)
	}
	router := internal.Router(database)

	log.Printf("starting server on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))

}
