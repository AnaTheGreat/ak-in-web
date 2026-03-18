package telemetry

import "log"

func Init(serviceName, env string) func() {
	log.Printf("telemetry initialized: service=%s env=%s", serviceName, env)

	// return shutdown hook
	return func() {
		log.Println("telemetry shutdown")
	}
}
