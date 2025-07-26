package integration

import (
	"net/http"
	"net/http/httptest"
	"personal_site/database"
	"personal_site/routers"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func setup(t *testing.T) {
	t.Setenv("DATABASE_DSN", ":memory:")
	t.Setenv("JWT_SECRET_KEY", "testsecretkey")
	t.Setenv("DEFAULT_TOKEN_EXPIRATION", "12h")

	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	router = gin.Default()
	routers.RegisterRouters(router, db)
}

func TestAuth(t *testing.T) {
	// Test the token generation and validation
	t.Run("Register, login, validate token and change-password", func(t *testing.T) {
		setup(t)
		// register
		w_reg := httptest.NewRecorder()
		req_reg, _ := http.NewRequest(http.MethodPost, "/auth/register",
			strings.NewReader(`{
				"email":"test@example.com", 
				"password":"password123",
				"nickname":"testuser"
			}`))

		router.ServeHTTP(w_reg, req_reg)

		assert.Equal(t, 200, w_reg.Code)

		// login
		w_login := httptest.NewRecorder()
		req_login, _ := http.NewRequest(http.MethodPost, "/auth/login",
			strings.NewReader(`{
				"email":"test@example.com",
				"password":"password123"
			}`))

		router.ServeHTTP(w_login, req_login)

		assert.Equal(t, 200, w_login.Code)

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

		// Validate token by change password
		w_change := httptest.NewRecorder()
		req_change, _ := http.NewRequest(http.MethodPost, "/auth/change-password",
			strings.NewReader(`{
				"old_password":"password123",
				"new_password":"newpassword123"
			}`))
		// Set the auth_token cookie instead of Authorization header
		req_change.AddCookie(authCookie)

		router.ServeHTTP(w_change, req_change)

		assert.Equal(t, 200, w_change.Code)

		// Check if the password was changed successfully
		w_oldPassword := httptest.NewRecorder()
		req_oldPassword, _ := http.NewRequest(http.MethodPost, "/auth/login",
			strings.NewReader(`{
				"email":"test@example.com",
				"password":"password123"
			}`))
		router.ServeHTTP(w_oldPassword, req_oldPassword)

		assert.Equal(t, 401, w_oldPassword.Code, "Old password should not be valid anymore")

		w_newPassword := httptest.NewRecorder()
		req_newPassword, _ := http.NewRequest(http.MethodPost, "/auth/login",
			strings.NewReader(`{
				"email":"test@example.com",
				"password":"newpassword123"
			}`))
		router.ServeHTTP(w_newPassword, req_newPassword)

		assert.Equal(t, 200, w_newPassword.Code, "New password should be valid")
	})
}
