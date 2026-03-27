package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	bibledb "github.com/sebdeveloper6952/bible-api/internal/db"
)

type bibleJSON struct {
	Version string               `json:"version"`
	Books   map[string]bookJSON  `json:"books"`
}

type bookJSON struct {
	Name     string                  `json:"name"`
	Number   int                     `json:"number"`
	Chapters map[string]chapterJSON  `json:"chapters"`
}

type chapterJSON struct {
	Number int                   `json:"number"`
	Verses map[string]verseJSON  `json:"verses"`
}

type verseJSON struct {
	Number  int    `json:"number"`
	Content string `json:"content"`
}

func main() {
	dbPath := "bible.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}
	jsonPath := "bible.json"
	if len(os.Args) > 2 {
		jsonPath = os.Args[2]
	}

	db, err := bibledb.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	f, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalf("open json: %v", err)
	}
	defer f.Close()

	var bible bibleJSON
	if err := json.NewDecoder(f).Decode(&bible); err != nil {
		log.Fatalf("decode json: %v", err)
	}

	if err := seed(db, bible); err != nil {
		log.Fatalf("seed: %v", err)
	}
	fmt.Println("Seeding complete.")
}

func seed(db *sql.DB, bible bibleJSON) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert version (idempotent)
	res, err := tx.Exec(
		`INSERT OR IGNORE INTO bible_versions (slug, name, language) VALUES (?, ?, ?)`,
		"rvr1960", bible.Version, "es",
	)
	if err != nil {
		return fmt.Errorf("insert version: %w", err)
	}
	versionID, err := res.LastInsertId()
	if err != nil || versionID == 0 {
		// already existed — fetch its id
		err = tx.QueryRow(`SELECT id FROM bible_versions WHERE slug = ?`, "rvr1960").Scan(&versionID)
		if err != nil {
			return fmt.Errorf("fetch version id: %w", err)
		}
	}

	// Sort book keys for deterministic insertion
	bookKeys := make([]string, 0, len(bible.Books))
	for k := range bible.Books {
		bookKeys = append(bookKeys, k)
	}
	sort.Slice(bookKeys, func(i, j int) bool {
		a, _ := strconv.Atoi(bookKeys[i])
		b, _ := strconv.Atoi(bookKeys[j])
		return a < b
	})

	for _, bk := range bookKeys {
		book := bible.Books[bk]
		res, err := tx.Exec(
			`INSERT OR IGNORE INTO books (version_id, number, name) VALUES (?, ?, ?)`,
			versionID, book.Number, book.Name,
		)
		if err != nil {
			return fmt.Errorf("insert book %d: %w", book.Number, err)
		}
		bookID, err := res.LastInsertId()
		if err != nil || bookID == 0 {
			err = tx.QueryRow(`SELECT id FROM books WHERE version_id = ? AND number = ?`, versionID, book.Number).Scan(&bookID)
			if err != nil {
				return fmt.Errorf("fetch book id: %w", err)
			}
		}

		// Sort chapter keys
		chapterKeys := make([]string, 0, len(book.Chapters))
		for k := range book.Chapters {
			chapterKeys = append(chapterKeys, k)
		}
		sort.Slice(chapterKeys, func(i, j int) bool {
			a, _ := strconv.Atoi(chapterKeys[i])
			b, _ := strconv.Atoi(chapterKeys[j])
			return a < b
		})

		for _, ck := range chapterKeys {
			chapter := book.Chapters[ck]
			res, err := tx.Exec(
				`INSERT OR IGNORE INTO chapters (book_id, number) VALUES (?, ?)`,
				bookID, chapter.Number,
			)
			if err != nil {
				return fmt.Errorf("insert chapter %d/%d: %w", book.Number, chapter.Number, err)
			}
			chapterID, err := res.LastInsertId()
			if err != nil || chapterID == 0 {
				err = tx.QueryRow(`SELECT id FROM chapters WHERE book_id = ? AND number = ?`, bookID, chapter.Number).Scan(&chapterID)
				if err != nil {
					return fmt.Errorf("fetch chapter id: %w", err)
				}
			}

			// Sort verse keys
			verseKeys := make([]string, 0, len(chapter.Verses))
			for k := range chapter.Verses {
				verseKeys = append(verseKeys, k)
			}
			sort.Slice(verseKeys, func(i, j int) bool {
				a, _ := strconv.Atoi(verseKeys[i])
				b, _ := strconv.Atoi(verseKeys[j])
				return a < b
			})

			for _, vk := range verseKeys {
				verse := chapter.Verses[vk]
				_, err := tx.Exec(
					`INSERT OR IGNORE INTO verses (chapter_id, number, content) VALUES (?, ?, ?)`,
					chapterID, verse.Number, verse.Content,
				)
				if err != nil {
					return fmt.Errorf("insert verse %d/%d/%d: %w", book.Number, chapter.Number, verse.Number, err)
				}
			}
		}
	}

	return tx.Commit()
}
