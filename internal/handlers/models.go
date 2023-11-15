package handlers

import (
	"time"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Apikey    string    `json:"apikey"`
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Apikey:    dbUser.Apikey,
	}
}

func databaseFeedToFeed(dbFeed database.Feed) Feed {
	return Feed{
		ID:        dbFeed.ID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
		Name:      dbFeed.Name,
		Url:       dbFeed.Url,
		UserID:    dbFeed.UserID,
	}
}

func databaseFeedsToFeeds(dbFeeds []database.Feed) []Feed {
	feeds := make([]Feed, len(dbFeeds))
	for i, dbFeed := range dbFeeds {
		feeds[i] = databaseFeedToFeed(dbFeed)
	}
	return feeds
}

func GetUserCreateParams(name string) database.CreateUserParams {
	return database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func GetFeedCreateParams(name, url string, userId uuid.UUID) database.CreateFeedParams {
	return database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    userId,
	}
}
