package routers

import (
	"personal_site/controllers"
	"personal_site/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router interface {
	RegisterRouters(r *gin.RouterGroup, db *gorm.DB)
}

func RegisterRouters(r *gin.Engine, db *gorm.DB) {
	var authRouter Router = auth{}
	authRouter.RegisterRouters(r.Group("/auth"), db)

	r.GET("/get-yt-data-api-token", middlewares.AuthOptional(), func(c *gin.Context) {
		controllers.GetYTDataAPIToken(c, db)
	})
}
