package token

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	tokenSign     = "titus.xapiens.id"
	exp       int = 1
)

type TokenData struct {
	Engine string
}

type JWTData struct {
	jwt.RegisteredClaims
	customClaims interface{}
}

// Init variable value
func initToken() {
	if os.Getenv("TOKEN_SIGN") != "" {
		tokenSign = os.Getenv("TOKEN_SIGN")
	}
	if os.Getenv("TOKEN_EXP") != "" {
		exp, _ = strconv.Atoi(os.Getenv("TOKEN_EXP"))
	}
}

func GenerateAccessToken(subject string, customClaims interface{}) (string, error) {
	initToken()

	tokenID, _ := uuid.NewUUID()
	claims := JWTData{
		RegisteredClaims: jwt.RegisteredClaims{
			// set token lifetime in timestamp
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(exp))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   subject,
			ID:        tokenID.String(),
		},
		customClaims: customClaims,
	}

	// generate a string using claims and HS256 algorithm
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the generated key using secretKey
	token, err := tokenString.SignedString([]byte(tokenSign))

	return token, err
}

func ValidateToken(token string) bool {
	initToken()

	claims := &JWTData{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSign), nil
	})

	if err != nil {
		log.Printf("token parse error %v", err)
		return false
	}

	log.Printf("logitimate access for ID %v, Subject %v", claims.ID, claims.Subject)
	return true
}

func ParseToken(token string) (*JWTData, error) {
	initToken()

	claims := &JWTData{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSign), nil
	})

	if err != nil {
		log.Printf("token parse error %v", err)
		return nil, err
	}

	log.Printf("logitimate access for ID %v, Subject %v", claims.ID, claims.Subject)
	return claims, err
}
