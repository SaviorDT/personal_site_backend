package controllers

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
	oauthgithub "golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

var githubOAuthConfig *oauth2.Config

func getGitHubOAuthConfig() (*oauth2.Config, error) {
	if githubOAuthConfig != nil {
		return githubOAuthConfig, nil
	}

	clientID, err := config.GetVariableAsString("GITHUB_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := config.GetVariableAsString("GITHUB_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	// Build redirect URL from shared path constant
	redirectURL := computeRedirectURL(apipaths.GitHubCallbackPath)
	githubOAuthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     oauthgithub.Endpoint,
		RedirectURL:  redirectURL,
	}

	return githubOAuthConfig, nil
}

// GitHubLoginStart redirects user to GitHub authorization URL
func GitHubLoginStart(c *gin.Context) {
	conf, err := getGitHubOAuthConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "GitHub OAuth not configured", "details": err.Error()})
		return
	}
	redirect := c.Query("redirect")
	state := encodeOAuthState(redirect)
	authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(302, authURL)
}

// GitHubLoginCallback handles GitHub redirect
func GitHubLoginCallback(c *gin.Context, db *gorm.DB) {
	conf, err := getGitHubOAuthConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "GitHub OAuth not configured", "details": err.Error()})
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

	ghUser, ghEmail, err := fetchGitHubUser(token.AccessToken)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch GitHub user", "details": err.Error()})
		return
	}

	// Ensure we have an email (GitHub can hide email)
	email := ghEmail
	if email == "" {
		// Create a synthetic email to satisfy NOT NULL + UNIQUE constraint
		email = fmt.Sprintf("github_%d@users.noreply.github.local", ghUser.ID)
	}

	user, err := ensureUserFromOAuth(db, models.AuthProviderGitHub, fmt.Sprintf("%d", ghUser.ID), email, ghUser.Login, ghUser.Name)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error", "details": err.Error()})
		return
	}

	// Generate our JWT token and set cookie
	jwtToken, err := GenerateToken(schemas.TokenPayload{UserID: user.ID, Role: string(user.Role), Nickname: user.Nickname}, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}
	exp, err := config.GetVariableAsTimeDuration("DEFAULT_TOKEN_EXPIRATION")
	if err != nil {
		exp = 12 * time.Hour
	}

	setAuthCookieForRedirect(c, jwtToken, exp, redirectBack)
	finalizeLoginResponse(c, redirectBack, user, "GitHub login successful")
}

type gitHubUser struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

type gitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func fetchGitHubUser(accessToken string) (*gitHubUser, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var u gitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, "", err
	}

	// Fetch emails
	req2, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req2.Header.Set("Authorization", "Bearer "+accessToken)
	req2.Header.Set("Accept", "application/vnd.github+json")
	resp2, err := client.Do(req2)
	if err != nil {
		return &u, "", nil // Ignore email error
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != 200 {
		return &u, "", nil
	}
	var emails []gitHubEmail
	if err := json.NewDecoder(resp2.Body).Decode(&emails); err == nil {
		for _, e := range emails {
			if e.Primary && e.Verified && e.Email != "" {
				return &u, e.Email, nil
			}
		}
		// fallback first
		if len(emails) > 0 && emails[0].Email != "" {
			return &u, emails[0].Email, nil
		}
	}
	return &u, "", nil
}
