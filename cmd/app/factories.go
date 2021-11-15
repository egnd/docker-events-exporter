package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	docker "github.com/docker/docker/client"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func NewCfg(path string, prefix string) (cfg *viper.Viper, err error) {
	cfg = viper.New()
	cfg.SetEnvPrefix(prefix)
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()
	cfg.SetConfigFile(path)
	err = cfg.ReadInConfig()

	return
}

func NewLogger(cfg *viper.Viper) *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.Local)
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().
		Logger().Level(zerolog.InfoLevel)

	if cfg.GetBool("debug") {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
			Level(zerolog.DebugLevel)
	}

	return &logger
}

func NewDockerClient(cfg *viper.Viper) (cli *docker.Client, err error) {
	cli, err = docker.NewClientWithOpts(
		docker.WithHost(cfg.GetString("docker.host")),
	)

	return
}

func NewHTTPServer(cfg *viper.Viper) (server *http.Server) {
	handler := http.NewServeMux()

	handler.Handle(cfg.GetString("metric.endpoint"), promhttp.Handler())

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.GetInt("port")),
		Handler: handler,
	}

	return
}
