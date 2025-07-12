package controllers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"fmt"
	"time"
)

var secretKey string

func GenerateToken(payload interface{}, id uint) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "deeelol_personal_site",
		"exp": jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
		"iat": jwt.NewNumericDate(now),
		"nbf": jwt.NewNumericDate(now),
		"aud": "deeelol_personal_site",
		"sub": id,
		"user": payload,
    })
    // Replace "secret" with your secret key

	key, err := getSecretKey()
	if err != nil {
		return "", fmt.Errorf("failed to get secret key: %v", err)
	}
    return token.SignedString(key)
}

func getSecretKey() ([]byte, error) {
	if secretKey == "" {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("Error loading .env file: %v", err)
		}

		secretKey := os.Getenv("JWT_SECRET_KEY")
		if secretKey == "" {
			return nil, fmt.Errorf("JWT_SECRET_KEY not set in .env")
		}
	}
	return []byte(secretKey), nil
}