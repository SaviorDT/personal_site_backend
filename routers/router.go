package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router interface {
	RegisterRouters(r *gin.RouterGroup, db *gorm.DB)
}

func RegisterRouters(r *gin.Engine, db *gorm.DB) {
	var authRouter Router = auth{}
	authRouter.RegisterRouters(r.Group("/auth"), db)
}
