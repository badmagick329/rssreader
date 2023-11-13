package handlers

import (
	"github.com/badmagick329/rssreader/internal/utils"
	"net/http"
)

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, 200, map[string]string{"status": "ok"})
}
