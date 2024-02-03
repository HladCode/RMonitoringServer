package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/HladCode/RMonitoringServer/internal/config"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/sensor/getPage"
	"github.com/HladCode/RMonitoringServer/internal/http-server/handlers/sensor/post"
	"github.com/HladCode/RMonitoringServer/internal/http-server/middleware/logger"
	"github.com/HladCode/RMonitoringServer/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	log.Info("starting refrigerator monitoring", slog.String("env", cnfg.Env))
	log.Debug("debug messages are enabled")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Post("/", post.New(log))
	router.Get("/idi", getPage.New(log, "/home/sergey/Projects/RMonitoringServer/internal/http-server/handlers/sensor/getPage/page/"))

	//	router.Get("/{alias}", rediect.New(log, storage))

	log.Info("starting server", slog.String("address", cnfg.Address))

	srv := &http.Server{
		Addr:         cnfg.Address,
		Handler:      router,
		ReadTimeout:  cnfg.HTTPServer.Timeout,
		WriteTimeout: cnfg.HTTPServer.Timeout,
		IdleTimeout:  cnfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", sl.Err(err))
	}

	log.Error("server stoped")
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
