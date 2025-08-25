package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"personal_site/config"
	"personal_site/controllers/storage"
	"personal_site/database"
	"personal_site/routers"
)

func main() {
	err := config.Init()
	if err != nil {
		panic(err)
	}

	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	r := gin.Default()

	startSetup()

	// CORS 配置
	allowedOrigins, _ := config.GetVariableAsString("CORS_ALLOWED_ORIGINS")

	corsConfig := cors.Config{
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"HEAD",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"X-Requested-With",
			"Cache-Control",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	if allowedOrigins != "" {
		corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
	}

	r.Use(cors.New(corsConfig))

	routers.RegisterRouters(r, db)

	r.Run(":80") // Start the server on port 8080
}

func startSetup() {
	// 清理 tmp 目錄
	tmpStoragePath, err := storage.GetStorageRoot()
	if err != nil {
		panic(err)
	}
	tmpStoragePath = filepath.Join(tmpStoragePath, "tmp")
	os.RemoveAll(tmpStoragePath)
}
