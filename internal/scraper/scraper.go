package scraper

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/google/uuid"
)

type Scraper struct {
	client              http.Client
	db                  *database.Queries
	concurrency         int32
	timeBetweetRequests time.Duration
}

func New(db *database.Queries) Scraper {
	return Scraper{
		client: http.Client{
			Timeout: time.Second * 10,
		},
		db:                  db,
		concurrency:         3,
		timeBetweetRequests: time.Second * 60,
	}
}

func (s *Scraper) makeRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Scraper) ReadFeed(url string) (*RSSFeed, error) {
	data, err := s.makeRequest(url)
	if err != nil {
		return nil, err
	}
	rssFeed, err := ParseRSSFeed(data)
	if err != nil {
		return nil, err
	}
	return rssFeed, nil
}

func (s *Scraper) ReadFeedFromFile(filename string) (*RSSFeed, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	bufReader := bytes.NewReader(data)
	decoder := xml.NewDecoder(bufReader)
	var rssFeed RSSFeed
	err = decoder.Decode(&rssFeed)
	if err != nil {
		return nil, err
	}
	return &rssFeed, nil
}

func (s *Scraper) StartScraping() {
	log.Printf("Starting scraping")
	ticker := time.NewTicker(s.timeBetweetRequests)
	wg := sync.WaitGroup{}
	for ; ; <-ticker.C {
		s.scrape(&wg)
	}
}

func (s *Scraper) scrape(wg *sync.WaitGroup) {
	ctx := context.Background()
	feeds, err := s.db.GetNextFeedsToFetch(ctx, s.concurrency)
	if err != nil {
		log.Printf("Error getting feeds %v", err)
		return
	}
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()
			rssFeed, err := s.ReadFeed(feed.Url)
			if err != nil {
				log.Printf("Error reading feed %v", err)
				return
			}
			for _, item := range rssFeed.Channel.Items {
				pubTime, err := time.Parse(time.RFC1123Z, item.PubDate)
				if err != nil {
					log.Printf("Error parsing time %v", err)
					continue
				}
				createParams := database.CreatePostParams{
					ID:          uuid.New(),
					FeedID:      feed.ID,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					Title:       item.Title,
					Url:         item.Link,
					Description: item.Description,
					PublishedAt: pubTime,
				}
				post, err := s.db.CreatePost(ctx, createParams)
				if err != nil {
					if strings.Contains(err.Error(), "duplicate key") {
						continue
					}
					log.Printf("Error creating post %v", err)
					continue
				}
				log.Printf("Created post %v, %s", post.ID, post.Title)
			}
		}(feed)
	}
	wg.Wait()
	log.Printf("Done scraping. Read %d feeds", len(feeds))
}
