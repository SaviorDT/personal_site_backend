package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"personal_site/config"
	"personal_site/database"
	"personal_site/routers"
	"personal_site/tasks"
)

func main() {
	// load env
	err := config.Init()
	if err != nil {
		panic(err)
	}

	// connect database
	db, err := database.InitDB()
	if err != nil {
		time.Sleep(30 * time.Second)
		db, err = database.InitDB()
		if err != nil {
			panic(err)
		}
	}

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
		MaxAge:           30 * 24 * time.Hour,
	}

	if allowedOrigins != "" {
		corsConfig.AllowOrigins = strings.Split(allowedOrigins, ",")
	}

	r := gin.Default()
	r.Use(cors.New(corsConfig))    // cors
	routers.RegisterRouters(r, db) // endpoints

	r.Run(":80") // Start the server on port 8080
}

func startSetup() {
	// gin debug mode
	ginmode, err := config.GetVariableAsString("GIN_MODE")
	if err != nil {
		gin.SetMode(gin.DebugMode)
	} else if ginmode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 清理 storage tmp 目錄
	tasks.ClearTmpStorage()
	// tmpStoragePath, err := storage.GetStorageRoot()
	// if err != nil {
	// 	panic(err)
	// }
	// tmpStoragePath = filepath.Join(tmpStoragePath, "tmp")
	// os.RemoveAll(tmpStoragePath)

	// set timezone
	if timezone, err := config.GetVariableAsString("TIMEZONE"); err == nil {
		loc, err := time.LoadLocation(timezone)
		if err != nil {
			panic(err)
		}
		time.Local = loc
	}
}
