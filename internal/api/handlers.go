package api

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/sebdeveloper6952/bible-api/internal/repository"
)

// internalError logs the error server-side and returns a generic 500 to the client.
func (s *Server) internalError(w http.ResponseWriter, r *http.Request, err error) {
	s.log.Error("handler error",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("error", err.Error()),
	)
	WriteError(w, http.StatusInternalServerError, "internal error")
}

// --- Version handlers ---

func (s *Server) listVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := s.versions.List()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, versions)
}

func (s *Server) getVersion(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("version")
	v, err := s.versions.GetBySlug(slug)
	if errors.Is(err, repository.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "version not found")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, v)
}

// --- Book handlers ---

func (s *Server) listBooks(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	books, err := s.books.ListByVersion(version)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, books)
}

func (s *Server) getBook(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	bookNum, err := strconv.Atoi(r.PathValue("book"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book number")
		return
	}
	book, err := s.books.GetByVersionAndNumber(version, bookNum)
	if errors.Is(err, repository.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, book)
}

// --- Chapter handlers ---

func (s *Server) listChapters(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	bookNum, err := strconv.Atoi(r.PathValue("book"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book number")
		return
	}
	chapters, err := s.chapters.ListByBook(version, bookNum)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, chapters)
}

func (s *Server) getChapter(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	bookNum, err := strconv.Atoi(r.PathValue("book"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book number")
		return
	}
	chapterNum, err := strconv.Atoi(r.PathValue("chapter"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid chapter number")
		return
	}
	chapter, err := s.chapters.GetWithVerses(version, bookNum, chapterNum)
	if errors.Is(err, repository.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "chapter not found")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, chapter)
}

// --- Verse handlers ---

func (s *Server) listVerses(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	bookNum, err := strconv.Atoi(r.PathValue("book"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book number")
		return
	}
	chapterNum, err := strconv.Atoi(r.PathValue("chapter"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid chapter number")
		return
	}
	verses, err := s.verses.ListByChapter(version, bookNum, chapterNum)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, verses)
}

func (s *Server) getVerse(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	bookNum, err := strconv.Atoi(r.PathValue("book"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid book number")
		return
	}
	chapterNum, err := strconv.Atoi(r.PathValue("chapter"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid chapter number")
		return
	}
	verseNum, err := strconv.Atoi(r.PathValue("verse"))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "invalid verse number")
		return
	}
	verse, err := s.verses.GetByNumber(version, bookNum, chapterNum, verseNum)
	if errors.Is(err, repository.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "verse not found")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, verse)
}
