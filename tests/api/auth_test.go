package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"personal_site/controllers"
	"personal_site/database"
	"personal_site/models"
	"personal_site/routers"
	"personal_site/schemas"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var router *gin.Engine
var db *gorm.DB

func setup(t *testing.T) {
	t.Setenv("DATABASE_DSN", ":memory:")
	t.Setenv("JWT_SECRET_KEY", "testsecretkey")
	t.Setenv("DEFAULT_TOKEN_EXPIRATION", "12h")

	var err error
	db, err = database.InitDB()
	if err != nil {
		panic(err)
	}

	router = gin.Default()
	routers.RegisterRouters(router, db)
}

func TestAuth(t *testing.T) {
	t.Run("Register", func(t *testing.T) {
		setup(t)
		w_reg := httptest.NewRecorder()
		req_reg, _ := http.NewRequest(http.MethodPost, "/auth/register",
			strings.NewReader(`{
				"email":"test-register@example.com", 
				"password":"password123",
				"nickname":"testuser"
			}`))

		router.ServeHTTP(w_reg, req_reg)

		assert.Equal(t, 200, w_reg.Code)

		// Check if user is created
		var user models.User
		err := db.First(&user, "email =  ?", "test-register@example.com").Error
		assert.NoError(t, err, "User should be created")
		assert.Equal(t, "testuser", user.Nickname, "User nickname should match")
		assert.Equal(t, models.RoleUser, user.Role, "User role should be 'user'")
		assert.Equal(t, models.AuthProviderPassword, user.Provider, "User provider should be 'password'")
		err2 := bcrypt.CompareHashAndPassword([]byte(user.Identifier), []byte("password123"))
		assert.NoError(t, err2, "User password should match")
	})

	t.Run("Login", func(t *testing.T) {
		setup(t)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		db.Create(&models.User{
			Nickname:   "testuser",
			Role:       models.RoleUser,
			Provider:   models.AuthProviderPassword,
			Email:      "test-login@example.com",
			Identifier: string(hashedPassword),
		})

		w_login := httptest.NewRecorder()
		req_login, _ := http.NewRequest(http.MethodPost, "/auth/login",
			strings.NewReader(`{
				"email":"test-login@example.com", 
				"password":"password123"
			}`))

		router.ServeHTTP(w_login, req_login)

		assert.Equal(t, 200, w_login.Code)

		// Check response data
		loginBodyBytes := w_login.Body.Bytes()
		var data map[string]any
		json.Unmarshal(loginBodyBytes, &data)
		assert.Equal(t, 1.0, data["user_id"], "User ID should be 1")
		assert.Equal(t, "user", data["role"], "Role should be 'user'")
		assert.Equal(t, "testuser", data["nickname"], "Nickname should match")

		// Check that auth_token cookie is set
		cookies := w_login.Result().Cookies()
		var authCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "auth_token" {
				authCookie = cookie
				break
			}
		}
		assert.NotNil(t, authCookie, "auth_token cookie should be set")
		assert.NotEmpty(t, authCookie.Value, "Token should not be empty")
	})

	t.Run("Change Password", func(t *testing.T) {
		setup(t)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		db.Create(&models.User{
			Nickname:   "testuser",
			Role:       models.RoleUser,
			Provider:   models.AuthProviderPassword,
			Email:      "test-change-password@example.com",
			Identifier: string(hashedPassword),
		})

		fakeToken, _ := controllers.GenerateToken(schemas.TokenPayload{
			UserID:   1,
			Role:     "user",
			Nickname: "testuser",
		}, 1)

		w_change := httptest.NewRecorder()
		req_change, _ := http.NewRequest(http.MethodPost, "/auth/change-password",
			strings.NewReader(`{
				"old_password":"password123",
				"new_password":"newpassword123"
			}`))

		req_change.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: fakeToken,
		})

		router.ServeHTTP(w_change, req_change)

		assert.Equal(t, 200, w_change.Code)

		// Check if the password was changed successfully
		var user models.User
		err := db.First(&user, "email =  ?", "test-change-password@example.com").Error
		assert.NoError(t, err, "User should be created")
		assert.Equal(t, "testuser", user.Nickname, "User nickname should match")
		assert.Equal(t, models.RoleUser, user.Role, "User role should be 'user'")
		assert.Equal(t, models.AuthProviderPassword, user.Provider, "User provider should be 'password'")
		err2 := bcrypt.CompareHashAndPassword([]byte(user.Identifier), []byte("newpassword123"))
		assert.NoError(t, err2, "User password should match")
	})

	t.Run("logout", func(t *testing.T) {
		setup(t)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		db.Create(&models.User{
			Nickname:   "testuser",
			Role:       models.RoleUser,
			Provider:   models.AuthProviderPassword,
			Email:      "test-change-password@example.com",
			Identifier: string(hashedPassword),
		})

		fakeToken, _ := controllers.GenerateToken(schemas.TokenPayload{
			UserID:   1,
			Role:     "user",
			Nickname: "testuser",
		}, 1)

		w_logout := httptest.NewRecorder()
		req_logout, _ := http.NewRequest(http.MethodPost, "/auth/logout", nil)
		req_logout.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: fakeToken,
		})

		router.ServeHTTP(w_logout, req_logout)

		assert.Equal(t, 200, w_logout.Code)

		// Check if the user is logged out
		cookies := w_logout.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "auth_token" {
				assert.True(t, cookie.MaxAge == -1, "auth_token cookie should be cleared after logout")
			}
		}
	})
}
