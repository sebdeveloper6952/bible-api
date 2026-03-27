package api

import (
	"log/slog"
	"net/http"

	"github.com/sebdeveloper6952/bible-api/internal/repository"
)

type Server struct {
	log      *slog.Logger
	versions *repository.VersionRepo
	books    *repository.BookRepo
	chapters *repository.ChapterRepo
	verses   *repository.VerseRepo
}

func NewServer(
	logger *slog.Logger,
	versions *repository.VersionRepo,
	books *repository.BookRepo,
	chapters *repository.ChapterRepo,
	verses *repository.VerseRepo,
) *Server {
	return &Server{
		log:      logger,
		versions: versions,
		books:    books,
		chapters: chapters,
		verses:   verses,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /versions", s.listVersions)
	mux.HandleFunc("GET /versions/{version}", s.getVersion)
	mux.HandleFunc("GET /versions/{version}/books", s.listBooks)
	mux.HandleFunc("GET /versions/{version}/books/{book}", s.getBook)
	mux.HandleFunc("GET /versions/{version}/books/{book}/chapters", s.listChapters)
	mux.HandleFunc("GET /versions/{version}/books/{book}/chapters/{chapter}", s.getChapter)
	mux.HandleFunc("GET /versions/{version}/books/{book}/chapters/{chapter}/verses", s.listVerses)
	mux.HandleFunc("GET /versions/{version}/books/{book}/chapters/{chapter}/verses/{verse}", s.getVerse)

	return loggingMiddleware(s.log)(mux)
}
