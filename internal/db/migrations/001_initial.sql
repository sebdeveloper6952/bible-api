-- +migrate Up

CREATE TABLE bible_versions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    slug       TEXT NOT NULL UNIQUE,
    name       TEXT NOT NULL,
    language   TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE books (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    version_id INTEGER NOT NULL REFERENCES bible_versions(id),
    number     INTEGER NOT NULL,
    name       TEXT NOT NULL,
    UNIQUE(version_id, number)
);

CREATE TABLE chapters (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    book_id INTEGER NOT NULL REFERENCES books(id),
    number  INTEGER NOT NULL,
    UNIQUE(book_id, number)
);

CREATE TABLE verses (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    chapter_id INTEGER NOT NULL REFERENCES chapters(id),
    number     INTEGER NOT NULL,
    content    TEXT NOT NULL,
    UNIQUE(chapter_id, number)
);

-- +migrate Down

DROP TABLE verses;
DROP TABLE chapters;
DROP TABLE books;
DROP TABLE bible_versions;
