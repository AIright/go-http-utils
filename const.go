package go_http_utils

import (
	"os"
	"strconv"
	"time"
)

const (
	envPort          = "SERVICE_PORT"
	envReadinessPort = "READINESS_PORT"

	envPodName                    = "POD_NAME"                               // runtime metrics
	envRuntimeMetricsInterval     = "GO_RUNTIME_METRICS_COLLECTION_INTERVAL" // runtime metrics interval
	defaultRuntimeMetricsInterval = 30 * time.Second

	defaultPort          = 8080
	defaultReadinessPort = 8081

	readinessProbeEndpoint = "/_info"
)

func envInt(env string, def int) int {
	if p, err := strconv.Atoi(os.Getenv(env)); err == nil {
		return p
	}
	return def
}

func envDuration(env string, def time.Duration) time.Duration {
	if d, err := time.ParseDuration(os.Getenv(env)); err == nil {
		return d
	}
	return def
}
