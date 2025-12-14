package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"personal_site/config"
	"personal_site/models"
	"personal_site/schemas"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginResponse struct {
	UserID   uint   `json:"user_id"`
	Message  string `json:"message"`
	Role     string `json:"role"`
	Nickname string `json:"nickname"`
	// Token        string `json:"token"`
	// RefreshToken string `json:"refresh_token,omitempty"` // 可選的刷新 token
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func Register(c *gin.Context, db *gorm.DB) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "unable to hash password"})
		return
	}

	user := models.User{
		Nickname:   req.Nickname,
		Role:       models.RoleUser,
		Provider:   models.AuthProviderPassword,
		Email:      req.Email,
		Identifier: string(hashedPassword), // In a real application, you should hash the password
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(200, gin.H{"message": "User registered successfully", "user_id": user.ID})
}

func Login(c *gin.Context, db *gorm.DB) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Attempt to login
	var user models.User
	err1 := db.Select("ID", "Role", "Nickname", "Identifier").Where("email = ?", req.Email).First(&user).Error // Cannot find user
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Identifier), []byte(req.Password))                       // Password mismatch

	// Login failed
	if err1 != nil || err2 != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// login successful
	token, err := GenerateTokenManual(schemas.TokenPayload{
		UserID:   user.ID,
		Role:     string(user.Role),
		Nickname: user.Nickname,
	}, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	fmt.Println("[Login Debug] Generated Token:", token) // DEBUG LOG

	setAuthCookie(c, token)

	c.JSON(200, loginResponse{
		UserID:   user.ID,
		Message:  "Login successful",
		Role:     string(user.Role),
		Nickname: user.Nickname,
	})
}

func Logout(c *gin.Context) {
	// 清除 auth_token cookie
	removeAuthCookie(c)

	c.JSON(200, gin.H{"message": "Logged out successfully"})
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ChangePassword(c *gin.Context, db *gorm.DB) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")

	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	tokenUser, ok := user.(schemas.TokenUser)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	var dbUser models.User
	if err := db.First(&dbUser, tokenUser.ID).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to find user"})
		return
	}

	if dbUser.Provider != models.AuthProviderPassword {
		c.JSON(403, gin.H{"error": "Password change only allowed for password-based accounts"})
		return
	}

	if !checkPasswordHash(req.OldPassword, dbUser.Identifier) {
		c.JSON(403, gin.H{"error": "Old password is incorrect"})
		return
	}

	newHashedPassword, err := hashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash new password"})
		return
	}

	dbUser.Identifier = newHashedPassword
	if err := db.Save(&dbUser).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

// ================= Helpers used by multiple providers =================

func fallbackNickname(parts ...string) string {
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			return p
		}
	}
	return "github_user"
}

// randomState simple random (reuse existing randomString)
func randomState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "state"
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func setAuthCookie(c *gin.Context, token string) {
	exp, err := config.GetVariableAsTimeDuration("DEFAULT_TOKEN_EXPIRATION")
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get token expiration duration", "details": err.Error()})
		return
	}

	// Determine if secure cookie should be used (e.g., based on environment or request scheme)
	// For simplicity, assuming HTTPS in production and HTTP in development
	secure := c.Request.URL.Scheme == "https" || c.Request.Header.Get("X-Forwarded-Proto") == "https"

	fmt.Println("[Cookie Debug] Setting Cookie:", token) // DEBUG LOG

	c.SetCookie(
		"auth_token_v2",    // cookie name
		token,              // cookie value
		int(exp.Seconds()), // max age in seconds
		"/",                // path
		"",                 // domain (empty means current domain)
		secure,             // secure (動態設定：開發環境 false，生產環境 true)
		true,               // httpOnly
	)
}

func removeAuthCookie(c *gin.Context) {
	c.SetCookie(
		"auth_token_v2", // cookie name
		"",              // empty value
		-1,              // max age -1 (delete immediately)
		"/",             // path
		"",              // domain (empty means current domain)
		true,            // secure (set to true in production with HTTPS)
		true,            // httpOnly
	)
}
