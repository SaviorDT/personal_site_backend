package schemas

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strconv"
    "time"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Payload TokenPayload `json:"payload"`
}

type TokenPayload struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Nickname string `json:"nickname"`
}

type TokenUser struct {
	ID       uint
	Role     string
	Nickname string
}

func NewTokenClaims[T interface{ ~string | ~uint }](sub T) *TokenClaims {
	var subject string
    switch v := any(sub).(type) {
    case string:
        subject = v
    case uint:
        subject = strconv.FormatUint(uint64(v), 10)
    }

	now := time.Now()
	return &TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims {
			Issuer: "https://後端.夢.台灣",
			Subject: subject,
			Audience: []string{"https://夢.台灣", "https://後端.夢.台灣"},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 12)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt: jwt.NewNumericDate(now),
			ID: randomString(16),
		},
	}
}

func (t *TokenPayload) ExtractUser() TokenUser {
	return TokenUser{
		ID:       t.UserID,
		Role:     t.Role,
		Nickname: t.Nickname,
	}
}

func randomString(n int) string {
    return uuid.NewString()
}