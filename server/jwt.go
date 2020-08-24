package main

import (
	"log"
	"time"
	"encoding/base64"

	"github.com/gobuffalo/uuid"
	"github.com/dgrijalva/jwt-go"
)

func generateRefresh() (refreshBase64 string){
	refreshT := uuid.Must(uuid.NewV4())
	js := Marsh(refreshT)
	refreshBase64 = base64.RawStdEncoding.EncodeToString(js)
	return refreshBase64
}

func generateAccess(ID uuid.UUID) (newAccessT *jwt.Token) {
	secretKey := []byte("owls_are_not_what_they_seem")
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(30 * time.Minute).Unix(),
	}
	newAccessT = jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	_, err := newAccessT.SignedString(secretKey)
	if err != nil {
		log.Fatal("error creating new access token", err)
	}
	return newAccessT
}

