package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "github.com/sebdeveloper6952/bible-api/docs"
	"github.com/sebdeveloper6952/bible-api/internal/api"
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

	dbPath := "/data/bible.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}
	addr := ":8080"
	if len(os.Args) > 2 {
		addr = os.Args[2]
	}

	db, err := bibledb.Open(dbPath)
	if err != nil {
		logger.Error("open db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	server := api.NewServer(
		logger,
		db,
		repository.NewVersionRepo(db),
		repository.NewBookRepo(db),
		repository.NewChapterRepo(db),
		repository.NewVerseRepo(db),
	)

	logger.Info("listening", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		logger.Error("server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
