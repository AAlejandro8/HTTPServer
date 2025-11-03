package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AAelajndro8/HTTPServer/internal"
	"github.com/AAelajndro8/HTTPServer/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	param1 := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(param1)
	if err != nil {
		responseWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	type response struct {
		Chirp
	}
	dbChirp, err := cfg.db.GetChirpByID(r.Context(),chirpID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "unable to get record", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Chirp: Chirp{
			Id: dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body: dbChirp.Body,
			UserId: dbChirp.UserID,
		},
	})
}


func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "unable to get chirps", err)
	}
	chrips := []Chirp{}
	for _, chirp := range dbChirps {
		chrips = append(chrips, Chirp{
				Id: chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body: chirp.Body,
				UserId: chirp.UserID,
			}) 
	}
	respondWithJSON(w, http.StatusOK, chrips)
	}



func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	// incoming structure (params)
	type parameters struct {
		Body string `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	// return values
	type response struct {
		Chirp
	}
	// decode the body 
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, "unable to decode", err)
		return
	}
	if len(params.Body) > 140 {
		responseWithError(w,http.StatusBadRequest, "Chirp is too long", nil)
		return 
	}
	badWords := map[string]struct{} {
		"kerfuffle": {},
		"sharbert": {},
		"fornax": {},
	}
	// check JWT Token
	bearerToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error getting bearer token",err)
		return
	}
	validatedID, err := internal.ValidateJWT(bearerToken, cfg.jwtsecret)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, fmt.Sprintf("error in the method %w", err), err)
		return
	}

	params.Body = replaceBadWords(params.Body, badWords)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: validatedID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to create chirp"))
	}
	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserId: chirp.UserID,
		},
	})
}

func replaceBadWords(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		if _, ok := badWords[word]; ok{
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}