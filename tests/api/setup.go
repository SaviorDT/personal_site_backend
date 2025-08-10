package api

import (
	"personal_site/database"
	"personal_site/routers"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router *gin.Engine
var db *gorm.DB

func setup(t *testing.T) {
	t.Setenv("DATABASE_DSN", ":memory:")
	t.Setenv("JWT_SECRET_KEY", "testsecretkey")
	t.Setenv("DEFAULT_TOKEN_EXPIRATION", "12h")
	t.Setenv("YT_DATA_API_TOKEN", "YT_DATA_API_TOKEN")

	var err error
	db, err = database.InitDB()
	if err != nil {
		panic(err)
	}

	router = gin.Default()
	routers.RegisterRouters(router, db)
}
