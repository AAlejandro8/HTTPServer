package internal

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	tokenIssuer = "chirpy-access"
	refreshTokenLength = 32
)

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	checked, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return checked, nil
}

// make the JWT with the claims
func MakeJWT(
	userID uuid.UUID, 
	tokenSecret string, 
	expiresIn time.Duration,
	) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.RegisteredClaims{
		Issuer: tokenIssuer,
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),
	})
	
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error){
	// make empty claims which parsewithclaims unmarshals into it (populating it)
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString, 
		&claims, 
		func(t *jwt.Token) (any, error) {return []byte(tokenSecret), nil},
	) 
	if err != nil {
		return uuid.UUID{},  err
	}

	if !token.Valid{
		return uuid.UUID{}, fmt.Errorf("token not valid")
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return uuid.UUID{}, err
	}

	if issuer != tokenIssuer {
		return uuid.UUID{}, errors.New("invalid issuer")
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid userID: %w", err)
	}

	return userID, nil
}

// getting the token out of the auth header field and stripping it
func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("authorization empty in headers")
	}
	rawTokenString := strings.TrimPrefix(strings.TrimSpace(tokenString), "Bearer ")
	
	return rawTokenString, nil
}

// make a refresh token which is used to make a new access token
func MakeRefreshToken() (string, error) {
	key := make([]byte, refreshTokenLength)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	} 

	hexString := hex.EncodeToString(key)

	return hexString, nil
}