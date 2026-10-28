package model

type Author struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
}

type AuthorFilter struct {
	Name  string 
	Page  int    
	Limit int    
}