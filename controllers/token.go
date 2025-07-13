package controllers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"os"
	"fmt"

	"personal_site/schemas"
)

var secretKey string

func GenerateToken(payload schemas.TokenPayload, id uint) (string, error) {
	// now := time.Now()
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
	// 	Iss: "deeelol_personal_site",
	// 	Exp: *jwt.NewNumericDate(now.Add(time.Hour * 12)),
	// 	Iat: *jwt.NewNumericDate(now),
	// 	Nbf: *jwt.NewNumericDate(now),
	// 	Aud: "deeelol_personal_site",
	// 	Sub: id,
	// 	Payload: payload,
    // })
	claims := schemas.NewTokenClaims(id)
	claims.Payload = payload

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

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

func ValidateToken(tokenString string) (*jwt.Token, error) {
	key, err := getSecretKey()
	if err != nil {
		return nil, fmt.Errorf("failed to get secret key: %v", err)
	}

	token, err := jwt.ParseWithClaims(tokenString, &schemas.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return token, nil
}