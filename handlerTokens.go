package main

import (
	"net/http"
	"time"
	"github.com/AAelajndro8/HTTPServer/internal"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// response
	type response struct {
		Token string `json:"token"`
	}
	// refresh token from the header
	refreshTokenString, err := internal.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error getting the refresh token", err)
		return 
	}
	// get the user we want to refresh the token for 
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshTokenString)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "error getting the user by token", err)
		return
	}
	// make the user a new token 
	newToken, err := internal.MakeJWT(
		user.ID, 
		cfg.jwtsecret, 
		time.Hour,
	) 
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error maing a new token", err)
		return
	}
	// send the user the new token
	respondWithJSON(w, http.StatusOK, response{
		Token: newToken,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	// get the refresh token from header
	refreshTokenString, err := internal.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error getting the refresh token", err)
		return 
	}
	// revoke the token 
	err = cfg.db.RevokeToken(r.Context(), refreshTokenString)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "error revoking the token", err)
		return
	}
	
	w.WriteHeader(204)
}