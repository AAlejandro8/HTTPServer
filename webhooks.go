package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/AAelajndro8/HTTPServer/internal"
	"github.com/google/uuid"
)

func (cfg *apiConfig) webHookRed(w http.ResponseWriter, r *http.Request) {
	// params
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, "error decoing request params", err)
		return 
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// check if valid api key
	apiKey, err := internal.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "error getting api key", err)
		return 
	}
	if apiKey != os.Getenv("POLKA_KEY") {
		responseWithError(w, http.StatusUnauthorized, "error getting api key", err)
		return 
	}
	// check if user exists 
	err = cfg.db.GetUserById(r.Context(), params.Data.UserID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "couldnt find the user", err)
		return
	} 
	// upgrade the user 
	_, err = cfg.db.UpgradeUserRed(r.Context(), params.Data.UserID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "unable to upgrade user ", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}