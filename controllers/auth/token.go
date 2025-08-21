package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"personal_site/config"
	"personal_site/schemas"
)

var getSecretKey = func() ([]byte, error) {
	key, err := config.GetVariableAsByteArr("JWT_SECRET_KEY")
	if err != nil {
		return nil, fmt.Errorf("failed to get secret key: %v", err)
	}
	return key, nil
}

func GenerateToken(payload schemas.TokenPayload, id uint) (string, error) {
	claims := schemas.NewTokenClaims(id)
	claims.Payload = payload

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key, err := getSecretKey()
	if err != nil {
		return "", fmt.Errorf("failed to get secret key: %v", err)
	}
	return token.SignedString(key)
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
