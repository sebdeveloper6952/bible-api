package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/sebdeveloper6952/bible-api/docs"
	"github.com/sebdeveloper6952/bible-api/internal/api"
	"github.com/sebdeveloper6952/bible-api/internal/config"
	bibledb "github.com/sebdeveloper6952/bible-api/internal/db"
	"github.com/sebdeveloper6952/bible-api/internal/repository"
)

// @title          Bible API
// @version        1.0
// @description    REST API to serve Bible content in multiple versions and languages.
// @host           localhost:8080
// @BasePath       /
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfgPath := "config/config.yaml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Error("load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db, err := bibledb.Open(cfg.Server.DBPath)
	if err != nil {
		logger.Error("open db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	server := api.NewServer(
		logger,
		cfg,
		db,
		repository.NewVersionRepo(db),
		repository.NewBookRepo(db),
		repository.NewChapterRepo(db),
		repository.NewVerseRepo(db),
	)

	logger.Info("listening", slog.String("addr", cfg.Server.Addr))
	if err := http.ListenAndServe(cfg.Server.Addr, server.Handler()); err != nil {
		logger.Error("server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
