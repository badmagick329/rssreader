package cmd

import (
	"fmt"
	"os"

	"github.com/badmagick329/rssreader/internal/handlers"
	"github.com/badmagick329/rssreader/internal/router"
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
	r.Run()
}
