package repository

import (
	"database/sql"
	"errors"

	"github.com/sebdeveloper6952/bible-api/internal/models"
)

type VerseRepo struct {
	db *sql.DB
}

func NewVerseRepo(db *sql.DB) *VerseRepo {
	return &VerseRepo{db: db}
}

func (r *VerseRepo) ListByChapter(versionSlug string, bookNumber, chapterNumber int) ([]models.Verse, error) {
	rows, err := r.db.Query(`
		SELECT vs.id, vs.chapter_id, vs.number, vs.content
		FROM verses vs
		JOIN chapters c ON c.id = vs.chapter_id
		JOIN books b ON b.id = c.book_id
		JOIN bible_versions v ON v.id = b.version_id
		WHERE v.slug = ? AND b.number = ? AND c.number = ?
		ORDER BY vs.number`, versionSlug, bookNumber, chapterNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []models.Verse
	for rows.Next() {
		var vs models.Verse
		if err := rows.Scan(&vs.ID, &vs.ChapterID, &vs.Number, &vs.Content); err != nil {
			return nil, err
		}
		verses = append(verses, vs)
	}
	return verses, rows.Err()
}

func (r *VerseRepo) GetByNumber(versionSlug string, bookNumber, chapterNumber, verseNumber int) (*models.Verse, error) {
	var vs models.Verse
	err := r.db.QueryRow(`
		SELECT vs.id, vs.chapter_id, vs.number, vs.content
		FROM verses vs
		JOIN chapters c ON c.id = vs.chapter_id
		JOIN books b ON b.id = c.book_id
		JOIN bible_versions v ON v.id = b.version_id
		WHERE v.slug = ? AND b.number = ? AND c.number = ? AND vs.number = ?`,
		versionSlug, bookNumber, chapterNumber, verseNumber,
	).Scan(&vs.ID, &vs.ChapterID, &vs.Number, &vs.Content)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &vs, nil
}
