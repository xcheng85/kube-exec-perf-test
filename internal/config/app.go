package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
	"os"
	"time"
)

type (
	AppConfig struct {
		Web             WebConfig
		ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"60s"`
	}
)

func NewAppConfig() (cfg AppConfig, err error) {
	// load env from file
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}
	// process env var to AppConfig
	err = envconfig.Process("", &cfg)
	return
}