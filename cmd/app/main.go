package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/egnd/docker-events-exporter/internal/devents"
	log "github.com/rs/zerolog/log"
)

var appVersion = "debug"

func main() {
	showVersion := flag.Bool("version", false, "Show app version.")
	cfgPath := flag.String("config", "configs/app.yaml", "Configuration file path.")
	cfgPrefix := flag.String("env-prefix", "DEE", "Prefix for env variables.")
	flag.Parse()

	if *showVersion {
		fmt.Println(appVersion)
		return
	}

	cfg, err := NewCfg(*cfgPath, *cfgPrefix)
	if err != nil {
		log.Fatal().Err(err).Msg("init config")
	}

	logger := NewLogger(cfg)

	dockerCLI, err := NewDockerClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("init docker listener")
	}

	go func() {
		logger.Info().Int("port", cfg.GetInt("port")).Msg("docker events listening...")
		if err := devents.NewListener(context.Background(), dockerCLI, logger, cfg).Listen(); err != nil {
			log.Fatal().Err(err).Msg("docker events listening")
		}
	}()

	server := NewHTTPServer(cfg)
	logger.Info().Int("port", cfg.GetInt("port")).Msg("http listening...")
	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err).Msg("http serving")
	}
}
