package scraper

import (
	"time"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/google/uuid"
)

func GetPostCreateParams(
	feed database.Feed,
	title, url, description string,
	publishedAt time.Time,
) database.CreatePostParams {
	return database.CreatePostParams{
		ID:          uuid.New(),
		FeedID:      feed.ID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       title,
		Url:         url,
		Description: description,
		PublishedAt: publishedAt,
	}
}
