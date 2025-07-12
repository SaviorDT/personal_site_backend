package controllers

import (
	"github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"personal_site/models"
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
	UserID uint   `json:"user_id"`
	Message string `json:"message"`
	Role string `json:"role"`
	Nickname string `json:"nickname"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"` // 可選的刷新 token
}

type tokenPayload struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Nickname string `json:"nickname"`
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
	c.JSON(201, gin.H{"message": "User registered successfully", "user_id": user.ID})
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
	err2 := bcrypt.CompareHashAndPassword([]byte(user.Identifier), []byte(req.Password)) // Password mismatch

	// Login failed
	if err1 != nil || err2 != nil {
		c.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// login successful
	token, err := GenerateToken(tokenPayload{
		UserID:   user.ID,
		Role:     string(user.Role),
		Nickname: user.Nickname,
	}, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token", "details": err.Error()})
		return
	}

	c.JSON(200, loginResponse{
		UserID:   user.ID,
		Message:  "Login successful",
		Role:     string(user.Role),
		Nickname: user.Nickname,	
		Token: token, // Replace with actual token generation logic
	})
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