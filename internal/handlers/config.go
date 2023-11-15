package handlers

import (
	"context"
	"database/sql"
	"log"

	"github.com/badmagick329/rssreader/internal/database"
	_ "github.com/lib/pq"
)

type Config struct {
	conn *sql.DB
	DB   *database.Queries
}

func New(dbUrl string) Config {
	if dbUrl == "" {
		log.Fatal("DB_URL not found in env file")
	}
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error connecting to db")
	}
	q := database.New(conn)
	return Config{conn: conn, DB: q}
}

func (cfg *Config) Close() error {
	return cfg.conn.Close()
}

func (cfg *Config) ClearDB(ctx context.Context) {
	cfg.DB.RemoveAllUsers(ctx)
	cfg.DB.RemoveAllFeeds(ctx)
}

