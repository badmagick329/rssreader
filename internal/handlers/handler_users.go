package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/badmagick329/rssreader/internal/auth"
	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/utils"
)

func (cfg *Config) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// query, err := url.ParseQuery(r.URL.RawQuery)
	// if err != nil {
	// 	log.Printf("Error parsing query: %s", err)
	// 	return
	// }
	// log.Printf("Query: %s", query)
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
	createParams := GetUserCreateParams(params.Name)
	// log.Printf("Creating user: %s", createParams.Name)
	user, err := cfg.DB.CreateUser(r.Context(), createParams)
	// log.Printf("Created user: %s", user.Name)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		utils.RespondWithError(w, 500, "Error creating user")
	}
	returnUser := databaseUserToUser(user)
	utils.RespondWithJSON(w, 201, returnUser)
}

func (cfg *Config) HandlerGetUserAuthed(
	w http.ResponseWriter,
	r *http.Request,
	user database.User,
) {
	utils.RespondWithJSON(w, 200, databaseUserToUser(user))
}

func (cfg *Config) HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		utils.RespondWithError(w, 400, "Error getting API key")
		return
	}
	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		utils.RespondWithError(w, 403, "Invalid API key")
		return
	}
	returnUser := databaseUserToUser(user)
	utils.RespondWithJSON(w, 200, returnUser)
}
