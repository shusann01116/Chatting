package model

import (
	"time"
)

type User struct {
	ID           string    `json:"userID"`
	Name         string    `json:"name"`
	ProfilePhoto string    `json:"profilephoto"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Rooms        []*Room   `json:"rooms"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
