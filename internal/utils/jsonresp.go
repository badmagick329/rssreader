package utils

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		// log.Printf("Failed to marshal: %v", payload)
		w.WriteHeader(500)
		return
	}
	// log.Printf("Response(%d): %s",code, data)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		// log.Printf("Server Error (%d): %s", code, msg)
	}
	RespondWithJSON(w, code, map[string]string{"error": msg})
}
