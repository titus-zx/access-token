package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	tokenSign       = "titus.xapiens.id"
	exp       int64 = 1
)

type JWTData struct {
	jwt.RegisteredClaims
}

func NewTokenRequest() {
	if os.Getenv("TOKEN_SIGN") != "" {
		tokenSign = os.Getenv("TOKEN_SIGN")
	}
	if os.Getenv("TOKEN_EXP") != "" {
		expVar, _ := strconv.Atoi(os.Getenv("TOKEN_EXP"))
		exp = int64(expVar)
	}
}

func GenerateAccessToken(subject string) (string, error) {
	// prepare claims for token
	tokenID, _ := uuid.NewUUID()
	claims := JWTData{
		RegisteredClaims: jwt.RegisteredClaims{
			// set token lifetime in timestamp
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(exp))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   subject,
			ID:        tokenID.String(),
		},
	}

	// generate a string using claims and HS256 algorithm
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the generated key using secretKey
	token, err := tokenString.SignedString([]byte(tokenSign))

	return token, err
}

func ValidateToken(token string) bool {
	claims := &JWTData{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSign), nil
	})

	if err != nil {
		log.Fatalf("token parse error %v", err)
		return false
	}

	now := time.Now()
	tokenExp := claims.ExpiresAt.Time

	if now.After(tokenExp) {
		log.Fatalf("token expired %v after now %v", now, tokenExp)
		return false
	}

	log.Printf("logitimate access for ID %v, Subject %v", claims.ID, claims.Subject)

	return true
}
