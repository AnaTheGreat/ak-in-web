package main

import (
	"log"
	"net/http"

	"ak-in-web/internal"
	"ak-in-web/internal/config"
	"ak-in-web/internal/telemetry"
)

func main() {
	cfg := config.Load()

	shutdownTelemetry := telemetry.Init(cfg.AppName, cfg.Env)
	defer shutdownTelemetry()

	addr := ":" + cfg.Port
	log.Printf("starting server on %s", addr)

	if err := http.ListenAndServe(addr, internal.Router()); err != nil {
		log.Fatal(err)
	}
}
