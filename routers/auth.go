package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

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

	r.POST("/change-password", middlewares.AuthRequired(), func(c *gin.Context) {
		controllers.ChangePassword(c, db)
	})
}
