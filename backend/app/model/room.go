package model

import "time"

type Room struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	Participants []User    `json:"participants"`
	Messages     []Message `json:"messages"`
}
