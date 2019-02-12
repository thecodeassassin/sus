package main

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/blendle/zapdriver"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"

	webserver "github.com/mycujoo/go-chi-webserver"
	"github.com/mycujoo/go-chi-webserver/middleware"
	"github.com/mycujoo/go-service-template/pkg/handler"
)

type Config struct {
	Env      string `env:"ENV" envDefault:"development"`
	Addr     string `env:"HTTP_ADDR" envDefault:":3000"`
}

func GenerateLogger() *zap.Logger {
	logger, err := zapdriver.NewProduction()

	if err != nil {
		log.Fatalf("Can't initialize logger: %v", err)
	}

	return logger
}

func GenerateConfig(logger *zap.Logger) *Config {
	// Ignore error if .env is missing
	err := godotenv.Load()

	if err != nil && !os.IsNotExist(err) {
		logger.Fatal("Unable to load .env", zap.Error(err))
	}

	cfg := &Config{}

	// Parse for built-in types
	if err := env.Parse(cfg); err != nil {
		logger.Fatal("Unable to parse env vars", zap.Error(err))
	}

	return cfg
}

func SetupRouter(logger *zap.Logger, env string) *chi.Mux {
	router := webserver.SetupRouter(env)

	cacheSize := 1000 * 1024 * 1024
	cacheHandler := middleware.CreateCacheHandler(cacheSize)

	router.Route("/test", func(r chi.Router) {
		r.Use(cacheHandler.Middleware)

		r.With(cacheHandler.Cache(60)).Get("/", handler.Get)
	})

	return router
}

func Run(logger *zap.Logger, addr string, router *chi.Mux) {
	logger.Info(fmt.Sprintf("Server listening on %s", addr))

	if err := webserver.Listen(addr, router); err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}

func main() {
	logger := GenerateLogger()
	defer logger.Sync()

	cfg := GenerateConfig(logger)
	router := SetupRouter(logger, cfg.Env)

	Run(logger, cfg.Addr, router)
}