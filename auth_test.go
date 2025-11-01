package main

import (
	"testing"
	"time"

	"github.com/AAelajndro8/HTTPServer/internal"
	"github.com/google/uuid"
)

func TestMakeJWT_Basic(t *testing.T) {
    userID := uuid.New()
    secret := "test-secret"
    exp := time.Minute


    token, err := internal.MakeJWT(userID, secret, exp)
    if err != nil {
	t.Fatalf("MakeJWT returned error: %v", err)
    }
    if token == "" {
	t.Fatal("tokenString is empty!")
    }
}

func TestValidateJWT_RoundTrip(t *testing.T) {
    userID := uuid.New()
    secret := "same-secret"
    exp := time.Minute

    token, err := internal.MakeJWT(userID, secret, exp)
    if err != nil {
	t.Fatalf("MakeJWT returned error: %v", err)
    }
    if token == "" {
	t.Fatal("tokenString is empty!")
    }
    gotID, err := internal.ValidateJWT(token, secret)
    if err != nil {
	t.Fatalf("ValidateJWT Failed error: %v", err)
    }
    if gotID != userID {
	t.Fatal("ID's do not match!")
    }
}

func TestValidateJWT_Expired(t *testing.T) {
    userID := uuid.New()
    secret := "secret"
    exp := time.Millisecond


    token, err := internal.MakeJWT(userID, secret, exp)
    if err != nil {
	t.Fatalf("MakeJWT returned error: %v", err)
    }
    if token == "" {
	t.Fatal("tokenString is empty!")
    }
    time.Sleep(2 * time.Millisecond)
    _, err = internal.ValidateJWT(token, secret)
    if err == nil {
	t.Fatalf("Expired token: %v", err)
    }
}

func TestValidateJWT_WrongSecret(t *testing.T) {
    userID := uuid.New()
    exp := time.Minute


    token, err := internal.MakeJWT(userID, "right", exp)
    if err != nil {
	t.Fatalf("MakeJWT returned error: %v", err)
    }
    if token == "" {
	t.Fatal("tokenString is empty!")
    }
    _, err = internal.ValidateJWT(token, "wrong")
    if err == nil {
	t.Fatalf("Signatre invalid: %v", err)
    }
}