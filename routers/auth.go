package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"personal_site/apipaths"
	"personal_site/controllers"
	"personal_site/middlewares"
)

type auth struct{}

func (a auth) RegisterRouters(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		controllers.Register(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, db)
	})

	r.POST("/logout", func(c *gin.Context) {
		controllers.Logout(c)
	})

	r.POST("/change-password", middlewares.AuthRequired(), func(c *gin.Context) {
		controllers.ChangePassword(c, db)
	})

	// GitHub OAuth
	r.GET(apipaths.GitHubLoginRel, func(c *gin.Context) {
		controllers.GitHubLoginStart(c)
	})
	r.GET(apipaths.GitHubCallbackRel, func(c *gin.Context) {
		controllers.GitHubLoginCallback(c, db)
	})

	// Google OAuth
	r.GET(apipaths.GoogleLoginRel, func(c *gin.Context) {
		controllers.GoogleLoginStart(c)
	})
	r.GET(apipaths.GoogleCallbackRel, func(c *gin.Context) {
		controllers.GoogleLoginCallback(c, db)
	})
}
