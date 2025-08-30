package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"personal_site/config"
	"personal_site/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// encodeOAuthState packs redirect and a nonce into a base64url JSON string
func encodeOAuthState(redirect string) string {
	payload := map[string]string{
		"n": randomState(),
		"r": redirect,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		delete(payload, "r") // If we can't marshal, just remove nonce
		b, err = json.Marshal(payload)
		if err != nil {
			b = []byte(randomState())
		}
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// decodeOAuthStateRedirect extracts the redirect from state
func decodeOAuthStateRedirect(raw string) string {
	if raw == "" {
		return ""
	}
	if b, err := base64.RawURLEncoding.DecodeString(raw); err == nil {
		var m map[string]string
		if json.Unmarshal(b, &m) == nil {
			if r, ok := m["r"]; ok {
				return r
			}
		}
	}
	return ""
}

// finalizeLoginResponse redirects back with user info in query, or returns JSON when no redirect
func finalizeLoginResponse(c *gin.Context, redirectBack string, user models.User, message string) {
	if redirectBack != "" {
		u, _ := url.Parse(redirectBack)
		q := u.Query()
		q.Set("login", "success")
		q.Set("user_id", strconv.Itoa(int(user.ID)))
		q.Set("message", message)
		q.Set("role", string(user.Role))
		q.Set("nickname", user.Nickname)
		u.RawQuery = q.Encode()
		c.Redirect(302, u.String())
		return
	}
	c.JSON(200, gin.H{"message": message, "user_id": user.ID, "role": user.Role, "nickname": user.Nickname})
}

// ensureUserFromOAuth finds or creates a user from provider + providerID
func ensureUserFromOAuth(db *gorm.DB, provider models.AuthProvider, providerID, email string, nicknameCandidates ...string) (models.User, error) {
	var user models.User
	if err := db.Where("provider = ? AND identifier = ?", provider, providerID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			user = models.User{
				Nickname:   fallbackNickname(nicknameCandidates...),
				Role:       models.RoleUser,
				Provider:   provider,
				Email:      email,
				Identifier: providerID,
			}
			if err := db.Create(&user).Error; err != nil {
				return models.User{}, err
			}
		} else {
			return models.User{}, err
		}
	}
	return user, nil
}

// computeRedirectURL builds an absolute callback URL based on the current request's scheme/host and a router path
func computeRedirectURL(callbackPath string) string {
	// Read base URL from environment (.env), e.g. PUBLIC_BASE_URL=https://example.com
	base, err := config.GetVariableAsString("PUBLIC_BASE_URL")
	if err != nil {
		// If not set, return the path as-is (caller should ensure env is configured)
		if !strings.HasPrefix(callbackPath, "/") {
			return "/" + callbackPath
		}
		return callbackPath
	}
	// Normalize base and join
	base = strings.TrimRight(base, "/")
	if !strings.HasPrefix(callbackPath, "/") {
		callbackPath = "/" + callbackPath
	}
	// Validate base is a URL
	if _, parseErr := url.Parse(base); parseErr != nil {
		return base + callbackPath
	}
	return base + callbackPath
}
