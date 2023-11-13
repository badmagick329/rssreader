package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/utils"
	"github.com/google/uuid"
)

func (cfg *Config) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var params parameters
	decoder.Decode(&params)
	if params == (parameters{}) {
		utils.RespondWithError(w, 400, "Invalid body")
		return
	}
	createParams := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user, err := cfg.DB.CreateUser(r.Context(), createParams)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		utils.RespondWithError(w, 500, "Error creating user")
	}
	returnUser := parameters{
		Name: user.Name,
	}
	utils.RespondWithJSON(w, 200, returnUser)
}
