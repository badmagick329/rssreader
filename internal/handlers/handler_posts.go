package handlers

import (
	"log"
	"net/http"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/utils"
)

func (cfg *Config) HandlerGetPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	dbPosts, err := cfg.DB.GetPostsByUserID(r.Context(), user.ID)
	if err != nil {
		log.Printf("Error getting posts: %s", err)
		utils.RespondWithError(w, 500, "Error getting posts")
		return
	}
	posts := databasePostsToPosts(dbPosts)
	utils.RespondWithJSON(w, 200, posts)
}
