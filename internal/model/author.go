package model

type Author struct {
    ID      int    `json:"id"`
    Name    string `json:"name"`
    Country string `json:"country"`
}