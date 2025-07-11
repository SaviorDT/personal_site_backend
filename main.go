package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	if err := InitDB(); err != nil {
		panic(err)
	}

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		db, err := DB.DB()
		if err != nil {
			c.String(http.StatusInternalServerError, "DB error: %v", err)
			return
		}
		err = db.Ping()
		if err != nil {
			c.String(http.StatusInternalServerError, "DB ping failed: %v", err)
			return
		}
		c.String(200, "DB connection OK!")
	})
	r.Run(":80") // Start the server on port 8080
}
