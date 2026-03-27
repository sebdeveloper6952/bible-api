package repository

import (
	"database/sql"
	"errors"

	"github.com/sebdeveloper6952/bible-api/internal/models"
)

type ChapterRepo struct {
	db *sql.DB
}

func NewChapterRepo(db *sql.DB) *ChapterRepo {
	return &ChapterRepo{db: db}
}

func (r *ChapterRepo) ListByBook(versionSlug string, bookNumber int) ([]models.Chapter, error) {
	rows, err := r.db.Query(`
		SELECT c.id, c.book_id, c.number
		FROM chapters c
		JOIN books b ON b.id = c.book_id
		JOIN bible_versions v ON v.id = b.version_id
		WHERE v.slug = ? AND b.number = ?
		ORDER BY c.number`, versionSlug, bookNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []models.Chapter
	for rows.Next() {
		var c models.Chapter
		if err := rows.Scan(&c.ID, &c.BookID, &c.Number); err != nil {
			return nil, err
		}
		chapters = append(chapters, c)
	}
	return chapters, rows.Err()
}

func (r *ChapterRepo) GetWithVerses(versionSlug string, bookNumber, chapterNumber int) (*models.Chapter, error) {
	var c models.Chapter
	err := r.db.QueryRow(`
		SELECT c.id, c.book_id, c.number
		FROM chapters c
		JOIN books b ON b.id = c.book_id
		JOIN bible_versions v ON v.id = b.version_id
		WHERE v.slug = ? AND b.number = ? AND c.number = ?`,
		versionSlug, bookNumber, chapterNumber,
	).Scan(&c.ID, &c.BookID, &c.Number)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(`
		SELECT id, chapter_id, number, content
		FROM verses WHERE chapter_id = ? ORDER BY number`, c.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var v models.Verse
		if err := rows.Scan(&v.ID, &v.ChapterID, &v.Number, &v.Content); err != nil {
			return nil, err
		}
		c.Verses = append(c.Verses, v)
	}
	return &c, rows.Err()
}
