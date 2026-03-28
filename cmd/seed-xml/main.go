package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"

	bibledb "github.com/sebdeveloper6952/bible-api/internal/db"
)

// bookNames maps canonical book numbers (1-66) to Spanish names.
// Used as fallback when the XML does not carry book name attributes.
var bookNames = map[int]string{
	1:  "Génesis",
	2:  "Éxodo",
	3:  "Levítico",
	4:  "Números",
	5:  "Deuteronomio",
	6:  "Josué",
	7:  "Jueces",
	8:  "Rut",
	9:  "1 Samuel",
	10: "2 Samuel",
	11: "1 Reyes",
	12: "2 Reyes",
	13: "1 Crónicas",
	14: "2 Crónicas",
	15: "Esdras",
	16: "Nehemías",
	17: "Ester",
	18: "Job",
	19: "Salmos",
	20: "Proverbios",
	21: "Eclesiastés",
	22: "Cantares",
	23: "Isaías",
	24: "Jeremías",
	25: "Lamentaciones",
	26: "Ezequiel",
	27: "Daniel",
	28: "Oseas",
	29: "Joel",
	30: "Amós",
	31: "Abdías",
	32: "Jonás",
	33: "Miqueas",
	34: "Nahúm",
	35: "Habacuc",
	36: "Sofonías",
	37: "Hageo",
	38: "Zacarías",
	39: "Malaquías",
	40: "San Mateo",
	41: "San Marcos",
	42: "San Lucas",
	43: "San Juan",
	44: "Hechos",
	45: "Romanos",
	46: "1 Corintios",
	47: "2 Corintios",
	48: "Gálatas",
	49: "Efesios",
	50: "Filipenses",
	51: "Colosenses",
	52: "1 Tesalonicenses",
	53: "2 Tesalonicenses",
	54: "1 Timoteo",
	55: "2 Timoteo",
	56: "Tito",
	57: "Filemón",
	58: "Hebreos",
	59: "Santiago",
	60: "1 Pedro",
	61: "2 Pedro",
	62: "1 Juan",
	63: "2 Juan",
	64: "3 Juan",
	65: "Judas",
	66: "Apocalipsis",
}

// ---- XML model ----

type xmlBible struct {
	Translation string         `xml:"translation,attr"`
	Testaments  []xmlTestament `xml:"testament"`
}

type xmlTestament struct {
	Name  string    `xml:"name,attr"`
	Books []xmlBook `xml:"book"`
}

type xmlBook struct {
	Number   int          `xml:"number,attr"`
	Name     string       `xml:"name,attr"` // present in some variants; empty in Beblia
	Chapters []xmlChapter `xml:"chapter"`
}

type xmlChapter struct {
	Number int        `xml:"number,attr"`
	Verses []xmlVerse `xml:"verse"`
}

type xmlVerse struct {
	Number  int    `xml:"number,attr"`
	Content string `xml:",chardata"`
}

// ---- main ----

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	dbPath := "data/bible.db"
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}
	xmlPath := "input/bible.xml"
	if len(os.Args) > 2 {
		xmlPath = os.Args[2]
	}
	// slug used to identify the version in the database (e.g. "tla")
	versionSlug := "tla"
	if len(os.Args) > 3 {
		versionSlug = os.Args[3]
	}

	logger.Info("opening database", slog.String("path", dbPath))
	db, err := bibledb.Open(dbPath)
	if err != nil {
		logger.Error("open db", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("loading bible XML", slog.String("path", xmlPath))
	f, err := os.Open(xmlPath)
	if err != nil {
		logger.Error("open xml", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer f.Close()

	var bible xmlBible
	if err := xml.NewDecoder(f).Decode(&bible); err != nil {
		logger.Error("decode xml", slog.String("error", err.Error()))
		os.Exit(1)
	}

	totalBooks := 0
	for _, t := range bible.Testaments {
		totalBooks += len(t.Books)
	}
	logger.Info("bible data loaded",
		slog.String("translation", bible.Translation),
		slog.Int("books", totalBooks),
	)

	if err := seed(db, bible, versionSlug, logger); err != nil {
		logger.Error("seed failed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	logger.Info("seeding complete")
}

func seed(db *sql.DB, bible xmlBible, slug string, logger *slog.Logger) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	logger.Info("inserting version", slog.String("slug", slug), slog.String("translation", bible.Translation))
	res, err := tx.Exec(
		`INSERT OR IGNORE INTO bible_versions (slug, name, language) VALUES (?, ?, ?)`,
		slug, bible.Translation, "es",
	)
	if err != nil {
		return fmt.Errorf("insert version: %w", err)
	}
	versionID, err := res.LastInsertId()
	if err != nil || versionID == 0 {
		err = tx.QueryRow(`SELECT id FROM bible_versions WHERE slug = ?`, slug).Scan(&versionID)
		if err != nil {
			return fmt.Errorf("fetch version id: %w", err)
		}
	}

	for _, testament := range bible.Testaments {
		for _, book := range testament.Books {
			name := book.Name
			if name == "" {
				name = bookNames[book.Number]
			}

			logger.Info("seeding book",
				slog.Int("number", book.Number),
				slog.String("name", name),
				slog.Int("chapters", len(book.Chapters)),
			)

			res, err := tx.Exec(
				`INSERT OR IGNORE INTO books (version_id, number, name) VALUES (?, ?, ?)`,
				versionID, book.Number, name,
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

			for _, chapter := range book.Chapters {
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

				for _, verse := range chapter.Verses {
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
	}

	logger.Info("committing transaction")
	return tx.Commit()
}
