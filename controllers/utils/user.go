package utils

import (
	"errors"
	"personal_site/schemas"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) uint {
	userID, err := GetUserIDStrict(c)
	if err != nil {
		return 0
	}
	return userID
}

func GetUserIDStrict(c *gin.Context) (uint, error) {
	tu, err := GetTokenUser(c)
	if err != nil {
		return 0, err
	}
	return tu.ID, nil
}

func GetUserNickname(c *gin.Context) string {
	userNickname, err := GetUserNicknameStrict(c)
	if err != nil {
		return "anonymous"
	}
	return userNickname
}

func GetUserNicknameStrict(c *gin.Context) (string, error) {
	tu, err := GetTokenUser(c)
	if err != nil {
		return "", err
	}
	return tu.Nickname, nil
}

// GetTokenUser returns the TokenUser stored in the context (set by auth middleware)
func GetTokenUser(c *gin.Context) (schemas.TokenUser, error) {
	if user, exists := c.Get("user"); exists {
		if userInfo, ok := user.(schemas.TokenUser); ok {
			return userInfo, nil
		}
	}
	return schemas.TokenUser{}, errors.New("user not authenticated")
}

// IsAdminUser returns true when the current user has admin role.
func IsAdminUser(c *gin.Context) bool {
	if user, exists := c.Get("user"); exists {
		if userInfo, ok := user.(schemas.TokenUser); ok {
			return userInfo.Role == "admin"
		}
	}
	return false
}
