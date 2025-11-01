package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AAelajndro8/HTTPServer/internal"
	"github.com/AAelajndro8/HTTPServer/internal/database"
	"github.com/google/uuid"
)

type User struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	// incoming data
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
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
	hash, err := internal.HashPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "hash didnt work", err)
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	})
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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	// incoming 
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}
	// response
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, "error decoding", err)
		return
	}
	// get user
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "user not found", err)
		return
	}
	// check password hashes
	exists, err := internal.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "check password hash method error", err)
		return
	}
	// no bueno
	if !exists {
		responseWithError(w, http.StatusUnauthorized, "password doesnt match", err)
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		},
	})
}