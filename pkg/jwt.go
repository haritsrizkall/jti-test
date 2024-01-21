package pkg

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func GenerateToken(userID int) (string, error) {
	claims := TokenPayload{
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ExtractPayload(token string) (TokenPayload, error) {
	secret := os.Getenv("JWT_SECRET")
	tokenPayload := TokenPayload{}
	_, err := jwt.ParseWithClaims(token, &tokenPayload, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return tokenPayload, err
	}
	return tokenPayload, nil
}
