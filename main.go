package main

import (
	"context"
	"flag"
	"net/http"
	"time"

	"github.com/flukis/inboice/services/infrastructure/querier"
	accountRegistration "github.com/flukis/inboice/services/internals/account_registration"
	"github.com/flukis/inboice/services/utils/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	var configFileName string
	flag.StringVar(&configFileName, "c", "config.yml", "Config file name")
	flag.Parse()

	cfg := config.LoadConfig(configFileName)
	log.Debug().Any("config", cfg).Msg("config loaded")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBCfg.ConnStr())
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	// querier database
	accountQuerier := querier.NewAccount(pool)

	// services
	registrationSvc := accountRegistration.New(accountQuerier)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)

	// router
	accountRegistrationHandler := accountRegistration.NewHttpHandler(registrationSvc)
	accountRegistrationHandler.Route(r)

	// Run server instance.
	log.Info().Msg("starting up server...")
	server := &http.Server{
		Handler:      r,
		Addr:         cfg.Listen.Addr(),
		ReadTimeout:  time.Second * time.Duration(cfg.Listen.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.Listen.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Listen.IdleTimeout),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}
	log.Info().Msg("server Stopped")
}
