package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	battleCatController "personal_site/controllers/battle_cat"
)

type battleCatRouter struct{}

func (battleCatRouter) RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.GET("/levels", func(c *gin.Context) {
		battleCatController.FilterLevels(c, db)
	})
}
