package model

type Author struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
}

type AuthorFilter struct {
	Name  string // ?name=  (search, optional)
	Page  int    
	Limit int    
}