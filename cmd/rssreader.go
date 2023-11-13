package cmd

import (
	"fmt"
	"os"

	"github.com/badmagick329/rssreader/internal/router"
	"github.com/joho/godotenv"
)

func Execute() {
	godotenv.Load()
	port := os.Getenv("PORT")
	fmt.Println("PORT is set to", port)
	r := router.New(port)
	r.AddRoute("GET", "/v1", "/healthz", router.HandlerReadiness)
	r.Run()
}
