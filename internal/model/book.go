package model

type Book struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    ReleaseYear int    `json:"release_year"`
    AuthorID    int    `json:"author_id"`
}