package main

import (
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
	routers.RegisterRouters(r, db)

	r.Run(":80") // Start the server on port 8080
}
