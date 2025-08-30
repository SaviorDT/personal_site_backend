package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"personal_site/apipaths"
	"personal_site/config"
	"personal_site/models"
	"personal_site/schemas"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	oauthgoogle "golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var googleOAuthConfig *oauth2.Config

func getGoogleOAuthConfig() (*oauth2.Config, error) {
	if googleOAuthConfig != nil {
		return googleOAuthConfig, nil
	}
	clientID, err := config.GetVariableAsString("GOOGLE_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := config.GetVariableAsString("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	// Build redirect URL from shared path constant
	redirectURL := computeRedirectURL(apipaths.GoogleCallbackPath)
	googleOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     oauthgoogle.Endpoint,
		RedirectURL:  redirectURL,
	}

	return googleOAuthConfig, nil
}

// GoogleLoginStart begins Google OAuth by redirecting to auth URL with encoded state
func GoogleLoginStart(c *gin.Context) {
	conf, err := getGoogleOAuthConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "Google OAuth not configured", "details": err.Error()})
		return
	}
	redirect := c.Query("redirect")
	state := encodeOAuthState(redirect)
	authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)
	fmt.Println("Redirecting to Google OAuth URL:", authURL)
	c.Redirect(302, authURL)
}

// GoogleLoginCallback handles Google redirect
func GoogleLoginCallback(c *gin.Context, db *gorm.DB) {
	conf, err := getGoogleOAuthConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "Google OAuth not configured", "details": err.Error()})
		return
	}
	code := c.Query("code")
	if code == "" {
		c.JSON(400, gin.H{"error": "Missing code"})
		return
	}
	redirectBack := decodeOAuthStateRedirect(c.Query("state"))
	if redirectBack == "" {
		redirectBack = c.Query("redirect")
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(401, gin.H{"error": "Code exchange failed", "details": err.Error()})
		return
	}

	gu, err := fetchGoogleUser(token.AccessToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch Google user", "details": err.Error()})
		return
	}

	email := gu.Email
	if email == "" {
		email = fmt.Sprintf("google_%s@users.noreply.google.local", gu.Sub)
	}

	user, err := ensureUserFromOAuth(db, models.AuthProviderGoogle, gu.Sub, email, gu.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error", "details": err.Error(), "user": gu})
		return
	}

	// Generate JWT and set cookie
	jwtToken, err := GenerateToken(schemas.TokenPayload{UserID: user.ID, Role: string(user.Role), Nickname: user.Nickname}, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	setAuthCookie(c, jwtToken)
	finalizeLoginResponse(c, redirectBack, user, "Google login successful")
}

type googleUser struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
}

func fetchGoogleUser(accessToken string) (*googleUser, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v3/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var u googleUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
