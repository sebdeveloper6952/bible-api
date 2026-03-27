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

// @Summary      List Bible versions
// @Description  Returns all available Bible versions.
// @Tags         versions
// @Produce      json
// @Success      200  {object}  versionListResponse
// @Failure      500  {object}  errorResponse
// @Router       /versions [get]
func (s *Server) listVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := s.versions.List()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, versions)
}

// @Summary      Get a Bible version
// @Description  Returns a single Bible version by its slug.
// @Tags         versions
// @Produce      json
// @Param        version  path      string  true  "Version slug (e.g. rvr1960)"
// @Success      200      {object}  versionResponse
// @Failure      404      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /versions/{version} [get]
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

// @Summary      List books
// @Description  Returns all books for a given Bible version.
// @Tags         books
// @Produce      json
// @Param        version  path      string  true  "Version slug (e.g. rvr1960)"
// @Success      200      {object}  bookListResponse
// @Failure      500      {object}  errorResponse
// @Router       /versions/{version}/books [get]
func (s *Server) listBooks(w http.ResponseWriter, r *http.Request) {
	version := r.PathValue("version")
	books, err := s.books.ListByVersion(version)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	WriteJSON(w, http.StatusOK, books)
}

// @Summary      Get a book
// @Description  Returns a single book by its number within a Bible version.
// @Tags         books
// @Produce      json
// @Param        version  path      string  true  "Version slug (e.g. rvr1960)"
// @Param        book     path      int     true  "Book number (1–66)"
// @Success      200      {object}  bookResponse
// @Failure      400      {object}  errorResponse
// @Failure      404      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /versions/{version}/books/{book} [get]
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

// @Summary      List chapters
// @Description  Returns all chapters for a given book and Bible version.
// @Tags         chapters
// @Produce      json
// @Param        version  path      string  true  "Version slug (e.g. rvr1960)"
// @Param        book     path      int     true  "Book number (1–66)"
// @Success      200      {object}  chapterListResponse
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /versions/{version}/books/{book}/chapters [get]
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

// @Summary      Get a chapter
// @Description  Returns a chapter with all its verses.
// @Tags         chapters
// @Produce      json
// @Param        version  path      string  true  "Version slug (e.g. rvr1960)"
// @Param        book     path      int     true  "Book number (1–66)"
// @Param        chapter  path      int     true  "Chapter number"
// @Success      200      {object}  chapterResponse
// @Failure      400      {object}  errorResponse
// @Failure      404      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /versions/{version}/books/{book}/chapters/{chapter} [get]
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

// @Summary      List verses
// @Description  Returns all verses for a given chapter.
// @Tags         verses
// @Produce      json
// @Param        version  path      string  true  "Version slug (e.g. rvr1960)"
// @Param        book     path      int     true  "Book number (1–66)"
// @Param        chapter  path      int     true  "Chapter number"
// @Success      200      {object}  verseListResponse
// @Failure      400      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /versions/{version}/books/{book}/chapters/{chapter}/verses [get]
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

// @Summary      Get a verse
// @Description  Returns a single verse by its number.
// @Tags         verses
// @Produce      json
// @Param        version  path      string  true  "Version slug (e.g. rvr1960)"
// @Param        book     path      int     true  "Book number (1–66)"
// @Param        chapter  path      int     true  "Chapter number"
// @Param        verse    path      int     true  "Verse number"
// @Success      200      {object}  verseResponse
// @Failure      400      {object}  errorResponse
// @Failure      404      {object}  errorResponse
// @Failure      500      {object}  errorResponse
// @Router       /versions/{version}/books/{book}/chapters/{chapter}/verses/{verse} [get]
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
