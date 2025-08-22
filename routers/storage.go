package routers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	storageController "personal_site/controllers/storage"
	"personal_site/middlewares"
)

type storageRouter struct{}

func (s storageRouter) RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	r.Use(middlewares.AuthOptional())

	// folder
	r.POST("/folder/*folder_name", func(c *gin.Context) {
		storageController.CreateFolder(c)
	})
	r.GET("/folder/*folder_name", func(c *gin.Context) {
		storageController.ListFolder(c)
	})
	r.PATCH("/folder/*folder_name", func(c *gin.Context) {
		storageController.UpdateFolder(c)
	})
	r.DELETE("/folder/*folder_name", func(c *gin.Context) {
		storageController.DeleteFolder(c)
	})
}
