package repository

import (
	"database/sql"
	"errors"

	"github.com/sebdeveloper6952/bible-api/internal/models"
)

type BookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *BookRepo {
	return &BookRepo{db: db}
}

func (r *BookRepo) ListByVersion(versionSlug string) ([]models.Book, error) {
	rows, err := r.db.Query(`
		SELECT b.id, b.version_id, b.number, b.name
		FROM books b
		JOIN bible_versions v ON v.id = b.version_id
		WHERE v.slug = ?
		ORDER BY b.number`, versionSlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.VersionID, &b.Number, &b.Name); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

func (r *BookRepo) GetByVersionAndNumber(versionSlug string, number int) (*models.Book, error) {
	var b models.Book
	err := r.db.QueryRow(`
		SELECT b.id, b.version_id, b.number, b.name
		FROM books b
		JOIN bible_versions v ON v.id = b.version_id
		WHERE v.slug = ? AND b.number = ?`, versionSlug, number,
	).Scan(&b.ID, &b.VersionID, &b.Number, &b.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}
