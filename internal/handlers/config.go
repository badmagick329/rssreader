package handlers

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Config struct {
	conn *sql.DB
	DB   *database.Queries
}

func New(dbUrl string) Config {
	// dbUrl := "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
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

func (cfg *Config) ClearDB(ctx context.Context) error {
	return cfg.DB.RemoveAllUsers(ctx)
}

func GetCreateParams(name string) database.CreateUserParams {
	return database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
