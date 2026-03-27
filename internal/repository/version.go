package repository

import (
	"database/sql"
	"errors"

	"github.com/sebdeveloper6952/bible-api/internal/models"
)

var ErrNotFound = errors.New("not found")

type VersionRepo struct {
	db *sql.DB
}

func NewVersionRepo(db *sql.DB) *VersionRepo {
	return &VersionRepo{db: db}
}

func (r *VersionRepo) List() ([]models.Version, error) {
	rows, err := r.db.Query(`SELECT id, slug, name, language FROM bible_versions ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []models.Version
	for rows.Next() {
		var v models.Version
		if err := rows.Scan(&v.ID, &v.Slug, &v.Name, &v.Language); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, rows.Err()
}

func (r *VersionRepo) GetBySlug(slug string) (*models.Version, error) {
	var v models.Version
	err := r.db.QueryRow(
		`SELECT id, slug, name, language FROM bible_versions WHERE slug = ?`, slug,
	).Scan(&v.ID, &v.Slug, &v.Name, &v.Language)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &v, nil
}
