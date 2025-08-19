package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"personal_site/apipaths"
	authController "personal_site/controllers/auth"
	"personal_site/middlewares"
)

type auth struct{}

func (a auth) RegisterRouters(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		authController.Register(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		authController.Login(c, db)
	})

	r.POST("/logout", func(c *gin.Context) {
		authController.Logout(c)
	})

	r.POST("/change-password", middlewares.AuthRequired(), func(c *gin.Context) {
		authController.ChangePassword(c, db)
	})

	// GitHub OAuth
	r.GET(apipaths.GitHubLoginRel, func(c *gin.Context) {
		authController.GitHubLoginStart(c)
	})
	r.GET(apipaths.GitHubCallbackRel, func(c *gin.Context) {
		authController.GitHubLoginCallback(c, db)
	})

	// Google OAuth
	r.GET(apipaths.GoogleLoginRel, func(c *gin.Context) {
		authController.GoogleLoginStart(c)
	})
	r.GET(apipaths.GoogleCallbackRel, func(c *gin.Context) {
		authController.GoogleLoginCallback(c, db)
	})
}
