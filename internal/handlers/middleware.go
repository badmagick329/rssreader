package handlers

import (
	"net/http"

	"github.com/badmagick329/rssreader/internal/auth"
	"github.com/badmagick329/rssreader/internal/database"
	"github.com/badmagick329/rssreader/internal/utils"
)

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *Config) MiddlewareAuth(handler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			utils.RespondWithError(w, 400, "Error getting API key")
			return
		}
		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}
		handler(w, r, user)
	}
}
