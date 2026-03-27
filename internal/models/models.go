package models

type Version struct {
	ID       int64  `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

type Book struct {
	ID        int64  `json:"id"`
	VersionID int64  `json:"-"`
	Number    int    `json:"number"`
	Name      string `json:"name"`
}

type Chapter struct {
	ID     int64   `json:"id"`
	BookID int64   `json:"-"`
	Number int     `json:"number"`
	Verses []Verse `json:"verses,omitempty"`
}

type Verse struct {
	ID        int64  `json:"id"`
	ChapterID int64  `json:"-"`
	Number    int    `json:"number"`
	Content   string `json:"content"`
}
