package auth

import (
	"fmt"
	"personal_site/schemas"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	t.Run("successful token generation", func(t *testing.T) {
		getSecretKey = func() ([]byte, error) {
			return []byte("fake_secret_value"), nil
		}

		payload := schemas.TokenPayload{
			UserID:   123,
			Role:     "admin",
			Nickname: "testuser",
		}

		tokenStr, err := GenerateToken(payload, 123)
		assert.NoError(t, err)

		// Test if token valid
		parsedToken, err := jwt.ParseWithClaims(tokenStr, &schemas.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("fake_secret_value"), nil
		})
		require.NoError(t, err)
		assert.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(*schemas.TokenClaims)
		require.True(t, ok)
		// 驗證 claims 的內容符合預期
		assert.Equal(t, "123", claims.Subject)
		assert.Equal(t, payload, claims.Payload)
	})

	t.Run("Cannot get secret key", func(t *testing.T) {
		getSecretKey = func() ([]byte, error) {
			return nil, fmt.Errorf("failed to get secret key")
		}
		_, err := GenerateToken(schemas.TokenPayload{}, 123)
		assert.Error(t, err)
	})
}

func TestValidateToken(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		getSecretKey = func() ([]byte, error) {
			return []byte("fake_secret_value"), nil
		}

		payload := schemas.TokenPayload{
			UserID:   123,
			Role:     "user",
			Nickname: "testuser",
		}

		tokenStr, err := GenerateToken(payload, 123)
		assert.NoError(t, err)

		validatedToken, err := ValidateToken(tokenStr)
		assert.NoError(t, err)
		assert.Equal(t, payload, validatedToken.Claims.(*schemas.TokenClaims).Payload)
	})

	t.Run("invalid token", func(t *testing.T) {
		getSecretKey = func() ([]byte, error) {
			return []byte("fake_secret_value"), nil
		}

		token, err := ValidateToken("invalid_token")
		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("Cannot get secret key", func(t *testing.T) {

		// Valid token
		payload := schemas.TokenPayload{
			UserID:   123,
			Role:     "user",
			Nickname: "testuser",
		}

		getSecretKey = func() ([]byte, error) {
			return []byte("fake_secret_value"), nil
		}

		tokenStr, err := GenerateToken(payload, 123)
		assert.NoError(t, err)

		getSecretKey = func() ([]byte, error) {
			return nil, fmt.Errorf("failed to get secret key")
		}

		validatedToken, err := ValidateToken(tokenStr)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
	})

	t.Run("token expires", func(t *testing.T) {
		getSecretKey = func() ([]byte, error) {
			return []byte("fake_secret_value"), nil
		}

		payload := schemas.TokenPayload{
			UserID:   123,
			Role:     "user",
			Nickname: "testuser",
		}

		claims := schemas.NewTokenClaims("123")
		claims.Payload = payload
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Second)) // Set expiration in the past

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte("fake_secret_value"))
		assert.NoError(t, err)

		validatedToken, err := ValidateToken(tokenStr)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
	})

	t.Run("token in future", func(t *testing.T) {
		getSecretKey = func() ([]byte, error) {
			return []byte("fake_secret_value"), nil
		}

		payload := schemas.TokenPayload{
			UserID:   123,
			Role:     "user",
			Nickname: "testuser",
		}

		claims := schemas.NewTokenClaims("123")
		claims.Payload = payload
		claims.NotBefore = jwt.NewNumericDate(time.Now().Add(time.Second)) // Set not before in the future

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte("fake_secret_value"))
		assert.NoError(t, err)

		validatedToken, err := ValidateToken(tokenStr)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
	})

	t.Run("Wrong secret key", func(t *testing.T) {
		getSecretKey = func() ([]byte, error) {
			return []byte("wrong_secret_key"), nil
		}

		payload := schemas.TokenPayload{
			UserID:   123,
			Role:     "user",
			Nickname: "testuser",
		}

		claims := schemas.NewTokenClaims("123")
		claims.Payload = payload

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString([]byte("fake_secret_value"))
		assert.NoError(t, err)

		validatedToken, err := ValidateToken(tokenStr)
		assert.Error(t, err)
		assert.Nil(t, validatedToken)
	})
}
