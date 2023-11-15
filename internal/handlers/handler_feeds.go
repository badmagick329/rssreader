package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/utils"
)

func (cfg *Config) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var params parameters
	decoder.Decode(&params)
	if params == (parameters{}) {
		utils.RespondWithError(w, 400, "Invalid body")
		return
	}
	createParams := GetFeedCreateParams(params.Name, params.Url, user.ID)
	feed, err := cfg.DB.CreateFeed(r.Context(), createParams)
	if err != nil {
		log.Printf("Error creating feed: %s", err)
		utils.RespondWithError(w, 500, "Error creating feed")
	}
	feedResponse := databaseFeedToFeed(feed)
	utils.RespondWithJSON(w, 201, feedResponse)
}
