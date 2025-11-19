package observability

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/uptrace-go/uptrace"

	"gaspoll/internal/config"
)

// Setup configures OpenTelemetry exporters (Uptrace) when DSN is provided.
func Setup(cfg config.Config) func(context.Context) error {
	if cfg.UptraceDSN == "" {
		log.Warn().Msg("UPTRACE_DSN is empty; traces and metrics are disabled")
		return func(context.Context) error { return nil }
	}

	uptrace.ConfigureOpentelemetry(
		uptrace.WithDSN(cfg.UptraceDSN),
		uptrace.WithServiceName(cfg.AppName),
		uptrace.WithServiceVersion("0.1.0"),
		uptrace.WithDeploymentEnvironment(cfg.AppEnv),
	)

	log.Info().Msg("Uptrace telemetry enabled")
	return uptrace.Shutdown
}
