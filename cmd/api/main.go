package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sebdeveloper6952/bible-api/internal/api"
	bibledb "github.com/sebdeveloper6952/bible-api/internal/db"
	"github.com/sebdeveloper6952/bible-api/internal/repository"
)

func main() {
	dbPath := "bible.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}
	addr := ":8080"
	if len(os.Args) > 2 {
		addr = os.Args[2]
	}

	db, err := bibledb.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	server := api.NewServer(
		repository.NewVersionRepo(db),
		repository.NewBookRepo(db),
		repository.NewChapterRepo(db),
		repository.NewVerseRepo(db),
	)

	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, server.Handler()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
