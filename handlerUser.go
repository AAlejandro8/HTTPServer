package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// incoming data
	type parameters struct {
		Email string `json:"email"`
	}
	// reponse 
	type response struct {
		User
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, "unable to decode", err)
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "unable to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated,response{
		User: User{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		},
	})
}