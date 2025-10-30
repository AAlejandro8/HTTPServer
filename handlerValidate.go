package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	// incoming structure (params)
	type parameters struct {
		Body string `json:"body"`
	}
	// return values
	type returnVals struct {
		Cleanedbody string `json:"cleaned_body"`
	}
	// decode the body 
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, "could't decode errors", err)
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
	params.Body = replaceBadWords(params.Body, badWords)
	respondWithJSON(w, http.StatusOK, returnVals{
		Cleanedbody: params.Body,
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