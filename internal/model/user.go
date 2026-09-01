package model


// User struct represents a user in our library system

type User struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`              //User answers won't bring back 
	Role         string `json:"role"`           // can be "member" or "admin"
    }                                        