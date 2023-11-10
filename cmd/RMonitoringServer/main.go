package main

import (
	"flag"
	"log"
	"log/slog"
	"os"

	"github.com/HladCode/RMonitoringServer/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ConfigPath := flag.String(
		"ConfigPath",
		"",
		"",
	)

	flag.Parse()

	if *ConfigPath == "" {
		log.Fatal("ConfigPath is not specified")
	}

	cnfg := config.MustLoad(*ConfigPath)

	log := setupLogger(cnfg.Env)
	log.Info("starting url-shortner", slog.String("env", cnfg.Env))
	log.Debug("debug messages are enabled")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelDebug,
				}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				}),
		)
	}

	return log
}
