package cmd

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/badmagick329/rssreader/internal/handlers"
	"github.com/badmagick329/rssreader/internal/router"
	"github.com/badmagick329/rssreader/internal/scraper"
	"github.com/joho/godotenv"
)

func Execute() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")
	fmt.Println("PORT is set to", port)
	r := router.New(port)
	cfg := handlers.New(dbUrl)
	r.AddRoute("GET", "/v1", "/healthz", handlers.HandlerReadiness)
	r.AddRoute("POST", "/v1", "/users", cfg.HandlerCreateUser)
	r.AddRoute("GET", "/v1", "/users", cfg.MiddlewareAuth(cfg.HandlerGetUserAuthed))
	r.AddRoute("POST", "/v1", "/feeds", cfg.MiddlewareAuth(cfg.HandlerCreateFeed))
	r.AddRoute("GET", "/v1", "/feeds", cfg.HandlerGetFeeds)
	r.AddRoute("POST", "/v1", "/feed_follows", cfg.MiddlewareAuth(cfg.HandlerCreateFeedFollow))
	r.AddRoute("DELETE", "/v1", "/feed_follows", cfg.MiddlewareAuth(cfg.HandlerDeleteFeedFollow))
	r.AddRoute("GET", "/v1", "/posts", cfg.MiddlewareAuth(cfg.HandlerGetPosts))
	scr := scraper.New(cfg.DB)
	go scr.StartScraping()
	r.Run()
}

func RSSTest() {
	scr := scraper.New(nil)
	rssFeed, err := scr.ReadFeed(scraper.BlogURL)
	if err != nil {
		log.Fatalf("Error Reading Feed %v", err)
	}
	b := []byte{}
	output := bytes.NewBuffer(b)
	encoder := xml.NewEncoder(output)
	err = encoder.Encode(rssFeed)
	if err != nil {
		log.Fatalf("Error Encoding %v", err)
	}
	os.WriteFile(scraper.BlogXML, output.Bytes(), 0644)
}

func RSSLoadTest() {
	scr := scraper.New(nil)
	rssFeed, err := scr.ReadFeedFromFile(scraper.BlogXML)
	if err != nil {
		log.Fatalf("Error reading feed from file %v", err)
	}
	log.Printf("%v", rssFeed.Channel.Items[0].PubDate)
}
