package routers

import (
	"personal_site/controllers"
	"personal_site/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router interface {
	RegisterRoutes(r *gin.RouterGroup, db *gorm.DB)
}

func RegisterRouters(r *gin.Engine, db *gorm.DB) {
	var authRouterVal Router = authRouter{}
	authRouterVal.RegisterRoutes(r.Group("/auth"), db)

	var storageRouterVal Router = storageRouter{}
	storageRouterVal.RegisterRoutes(r.Group("/storage"), db)

	var battleCatRouterVal Router = battleCatRouter{}
	battleCatRouterVal.RegisterRoutes(r.Group("/battle-cat"), db)

	r.GET("/get-yt-data-api-token", middlewares.AuthOptional(), func(c *gin.Context) {
		controllers.GetYTDataAPIToken(c, db)
	})
}
