package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ServiceConfig struct {
	General    General
	HTTPServer HTTPServer
}

type General struct {
	Env      string `envconfig:"ENV" default:"local"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

type HTTPServer struct {
	Host                    string        `envconfig:"HTTP_HOST" default:"0.0.0.0"`
	Port                    string        `envconfig:"HTTP_PORT" default:"4080"`
	GracefulShutdownTimeout time.Duration `envconfig:"HTTP_GRACEFUL_SHUTDOWN_TIMEOUT" default:"15s"`
	WriteTimeout            time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"30s"`
	ReadTimeout             time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"30s"`
	ReadHeaderTimeout       time.Duration `envconfig:"HTTP_READ_HEADER_TIMEOUT" default:"30s"`
	IdleTimeout             time.Duration `envconfig:"HTTP_IDLE_TIMEOUT" default:"300s"`
	// Be careful with using the above timeouts. They do not cancel the context, hence the handler
	// execution will continiue and resources would not be released. They just set the deadline for
	// respective read and write operations.
	// In order to set the timeout for the handler, consider on setting the handler timeout and use
	// TimeoutMiddleware
	HandlerTimeout time.Duration `envconfig:"HTTP_HANDLER_TIMEOUT" default:"15s"`
}

func Load() (*ServiceConfig, error) {
	cfg := &ServiceConfig{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse configuration: %v", err)
	}

	return cfg, nil
}
