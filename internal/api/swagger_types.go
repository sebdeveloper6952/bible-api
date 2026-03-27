package api

import "github.com/sebdeveloper6952/bible-api/internal/models"

// Typed response wrappers for swaggo annotations.
// Not used at runtime — only referenced in doc comments.

type versionResponse     struct { Data *models.Version  `json:"data"` }
type versionListResponse struct { Data []models.Version `json:"data"` }
type bookResponse        struct { Data *models.Book     `json:"data"` }
type bookListResponse    struct { Data []models.Book    `json:"data"` }
type chapterResponse     struct { Data *models.Chapter  `json:"data"` }
type chapterListResponse struct { Data []models.Chapter `json:"data"` }
type verseResponse       struct { Data *models.Verse    `json:"data"` }
type verseListResponse   struct { Data []models.Verse   `json:"data"` }
type healthResponse      struct { Data map[string]string `json:"data"` }
type errorResponse       struct { Error string           `json:"error"` }
