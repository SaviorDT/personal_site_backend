package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"personal_site/controllers"
)

type auth struct{}

func (a auth) RegisterRouters(r *gin.RouterGroup, db *gorm.DB) {
	r.POST("/register", func(c *gin.Context) {
		controllers.Register(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		controllers.Login(c, db)
	})
}
