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
	var userID *uint
	if user, exists := c.Get("user"); exists {
		if userInfo, ok := user.(schemas.TokenUser); ok {
			userID = &userInfo.ID
		}
	}

	if userID == nil {
		return 0, errors.New("user not authenticated")
	}

	return *userID, nil
}

func GetUserNickname(c *gin.Context) string {
	userNickname, err := GetUserNicknameStrict(c)
	if err != nil {
		return "anonymous"
	}
	return userNickname
}

func GetUserNicknameStrict(c *gin.Context) (string, error) {
	var userNickname *string
	if user, exists := c.Get("user"); exists {
		if userInfo, ok := user.(schemas.TokenUser); ok {
			userNickname = &userInfo.Nickname
		}
	}

	if userNickname == nil {
		return "", errors.New("user not authenticated")
	}

	return *userNickname, nil
}
