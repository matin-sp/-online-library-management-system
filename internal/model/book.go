package model

type Book struct {
    ID          int    `json:"id"`
    Title       string `json:"title"`
    ReleaseYear int    `json:"release_year"`
    AuthorID    int    `json:"author_id"`
}


type BookFilter struct {
	Title    string 
	AuthorID int    
	YearFrom int   
	YearTo   int    
	Page     int 
	Limit    int    
}