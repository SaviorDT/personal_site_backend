package controllers

import (
	"net/http"
	"personal_site/config"
	"personal_site/models"
	"personal_site/schemas"
	"slices"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetTokenResponse struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func GetYTDataAPIToken(c *gin.Context, db *gorm.DB) {
	// Get parameters from query string
	q1 := c.Query("question1")
	q2 := c.Query("question2")
	q3 := c.Query("question3")

	// Validate that all required parameters are provided
	if q1 == "" || q2 == "" || q3 == "" {
		c.JSON(http.StatusBadRequest, GetTokenResponse{
			Token:   "",
			Success: false,
			Message: "All questions (question1, question2, question3) are required",
		})
		return
	}

	// Validate parameter lengths (should not exceed 256 characters)
	if len(q1) > 256 || len(q2) > 256 || len(q3) > 256 {
		c.JSON(http.StatusBadRequest, GetTokenResponse{
			Token:   "",
			Success: false,
			Message: "Question answers must not exceed 256 characters",
		})
		return
	}

	// Get user information if logged in
	var userID *uint
	if user, exists := c.Get("user"); exists {
		if userInfo, ok := user.(schemas.TokenUser); ok {
			userID = &userInfo.ID
		}
	}

	history := models.YTDataAPITokenHistory{
		Q1:     q1,
		Q2:     q2,
		Q3:     q3,
		UserID: userID,
	}

	if err := db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, GetTokenResponse{
			Token:   "",
			Success: false,
			Message: "Failed to save token history",
		})
		return
	}

	// Call the YouTube Data API to get the token

	// Define valid answers for Q2
	validAnswers := []string{"Savior_DT", "Savior_TD", "literal_sorcerer", "秋滲幽人", "冬滲冬瓜", "開學忙死了", "tony20040424"}

	// Check if q2 matches one of the valid answers
	isValidQ2 := slices.Contains(validAnswers, q2)

	// If Q2 is not valid, return error
	if !isValidQ2 {
		c.JSON(http.StatusBadRequest, GetTokenResponse{
			Token:   "",
			Success: false,
			Message: "答案錯誤",
		})
		return
	}

	token, err := config.GetVariableAsString("YT_DATA_API_TOKEN")
	if err != nil || token == "" {
		c.JSON(http.StatusInternalServerError, GetTokenResponse{
			Token:   "",
			Success: false,
			Message: "Failed to get token",
		})
		return
	}

	c.JSON(http.StatusOK, GetTokenResponse{
		Token:   token,
		Success: true,
		Message: "Token retrieved successfully",
	})
}
