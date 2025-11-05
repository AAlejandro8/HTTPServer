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
	Is_chirpy_red bool `json:"is_chirpy_red"`
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
		Is_chirpy_red: user.IsChirpyRed,
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
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, "error decoding request params", err)
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
		return
	}

	// time the token lasts 
	expirationTime := time.Hour
	// make the token 
	accessToken, err := internal.MakeJWT(user.ID, cfg.jwtsecret, expirationTime)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error making JWT token", err)
		return
	}
	// make refresh token
	refreshString, err := internal.MakeRefreshToken() 
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error making refresh token", err)
		return
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshString,
		UserID: user.ID,
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error storing the refresh token", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
		Id: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Is_chirpy_red: user.IsChirpyRed,
		},
		Token: accessToken,
		RefreshToken: refreshString,
	})
}

func (cfg *apiConfig) handlerUpdateInfo (w http.ResponseWriter, r *http.Request) {
	// incoming 
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	// response 
	type response struct {
		User
	}
	// check access token in header
	accessToken, err := internal.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "no access token in request header", err)
		return
	}
	// validate the token 
	userID, err := internal.ValidateJWT(accessToken, cfg.jwtsecret)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "token isnt valid", err)
		return
	}

	// get the request params
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		responseWithError(w, http.StatusInternalServerError, "error decoing request params", err)
		return 
	}

	// hash new password
	hashedPassword, err := internal.HashPassword(params.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error hashing the password", err)
		return
	}

	// update the email and password 
	updatedUser, err := cfg.db.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "error updating user info", err)
		return 
	}
	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id: updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email: updatedUser.Email,
			Is_chirpy_red: updatedUser.IsChirpyRed,
		},
	})
}