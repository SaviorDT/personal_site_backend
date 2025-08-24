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
	r.POST("/folder/*folder_path", func(c *gin.Context) {
		storageController.CreateFolder(c)
	})
	r.GET("/folder/*folder_path", func(c *gin.Context) {
		storageController.ListFolder(c)
	})
	r.PATCH("/folder/*folder_path", func(c *gin.Context) {
		storageController.UpdateFolder(c)
	})
	r.DELETE("/folder/*folder_path", func(c *gin.Context) {
		storageController.DeleteFolder(c)
	})

	// file
	r.GET("/file/*file_path", func(c *gin.Context) {
		storageController.GetFile(c)
	})
	r.POST("/file/*file_path", func(c *gin.Context) {
		storageController.UploadFile(c)
	})
	r.PATCH("/file/*file_path", func(c *gin.Context) {
		storageController.UpdateFile(c)
	})
	r.DELETE("/file/*file_path", func(c *gin.Context) {
		storageController.DeleteFile(c)
	})
}
