package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"personal_site/config"
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

	// CORS 配置
	allowedOrigins, _ := config.GetVariableAsString("CORS_ALLOWED_ORIGINS")

	corsConfig := cors.Config{
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
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
