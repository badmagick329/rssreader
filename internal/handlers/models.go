package handlers

import (
	"time"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Apikey     string    `json:"apikey"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:         dbUser.ID,
		Name:       dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Apikey:     dbUser.Apikey,
	}
}
