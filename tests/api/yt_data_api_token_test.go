package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	authController "personal_site/controllers/auth"
	"personal_site/models"
	"personal_site/schemas"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetYTDataAPIToken(t *testing.T) {
	t.Run("Get YT Data API Token", func(t *testing.T) {
		setup(t)

		// Build URL with query parameters
		params := url.Values{}
		params.Set("question1", "test user")
		params.Set("question2", "Savior_DT")
		params.Set("question3", "test")

		testURL := "/get-yt-data-api-token?" + params.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testURL, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		// Test token right.
		loginBodyBytes := w.Body.Bytes()
		var data map[string]any
		json.Unmarshal(loginBodyBytes, &data)
		assert.Equal(t, "YT_DATA_API_TOKEN", data["token"], "Token should match")

		// Check database history.
		var history models.YTDataAPITokenHistory
		err := db.First(&history).Error
		assert.NoError(t, err, "History should be created")
		assert.Equal(t, "test user", history.Q1, "Q1 should match")
		assert.Equal(t, "Savior_DT", history.Q2, "Q2 should match")
		assert.Equal(t, "test", history.Q3, "Q3 should match")
		assert.Equal(t, uint(0), *history.UserID, "UserID should be nil for unauthenticated request")
	})

	t.Run("Get YT Data API Token with invalid Q2", func(t *testing.T) {
		setup(t)

		// Build URL with query parameters
		params := url.Values{}
		params.Set("question1", "test user")
		params.Set("question2", "invalid_answer")
		params.Set("question3", "test")

		testURL := "/get-yt-data-api-token?" + params.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testURL, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, 400, w.Code)

		loginBodyBytes := w.Body.Bytes()
		var data map[string]any
		json.Unmarshal(loginBodyBytes, &data)
		assert.Equal(t, "", data["token"], "Token should be empty")
		assert.Equal(t, false, data["success"], "Success should be false")

		// Check database history.
		var history models.YTDataAPITokenHistory
		err := db.First(&history).Error
		assert.NoError(t, err, "History should be created")
		assert.Equal(t, "test user", history.Q1, "Q1 should match")
		assert.Equal(t, "invalid_answer", history.Q2, "Q2 should match")
		assert.Equal(t, "test", history.Q3, "Q3 should match")
		assert.Equal(t, uint(0), *history.UserID, "UserID should be nil for unauthenticated request")
	})

	t.Run("Get YT Data API Token with authenticated user", func(t *testing.T) {
		setup(t)

		// Create a test user
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		testUser := models.User{
			Nickname:   "testuser",
			Role:       models.RoleUser,
			Provider:   models.AuthProviderPassword,
			Email:      "test-token@example.com",
			Identifier: string(hashedPassword),
		}
		db.Create(&testUser)

		// Generate token for the test user
		fakeToken, _ := authController.GenerateToken(schemas.TokenPayload{
			UserID:   testUser.ID,
			Role:     string(testUser.Role),
			Nickname: testUser.Nickname,
		}, testUser.ID)

		// Build URL with query parameters
		params := url.Values{}
		params.Set("question1", "test user")
		params.Set("question2", "Savior_DT")
		params.Set("question3", "test")

		testURL := "/get-yt-data-api-token?" + params.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testURL, nil)

		// Add authentication cookie
		req.AddCookie(&http.Cookie{
			Name:  "auth_token",
			Value: fakeToken,
		})

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		// Check database history with user ID
		var history models.YTDataAPITokenHistory
		err := db.First(&history).Error
		assert.NoError(t, err, "History should be created")
		assert.Equal(t, "test user", history.Q1, "Q1 should match")
		assert.Equal(t, "Savior_DT", history.Q2, "Q2 should match")
		assert.Equal(t, "test", history.Q3, "Q3 should match")
		assert.NotNil(t, history.UserID, "UserID should not be nil for authenticated request")
		assert.Equal(t, testUser.ID, *history.UserID, "UserID should match the authenticated user")
	})

	t.Run("Get YT Data API Token without authentication", func(t *testing.T) {
		setup(t)

		// Build URL with query parameters
		params := url.Values{}
		params.Set("question1", "test user")
		params.Set("question2", "Savior_DT")
		params.Set("question3", "test")

		testURL := "/get-yt-data-api-token?" + params.Encode()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testURL, nil)
		// No authentication cookie

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		// Check database history without user ID
		var history models.YTDataAPITokenHistory
		err := db.First(&history).Error
		assert.NoError(t, err, "History should be created")
		assert.Equal(t, "test user", history.Q1, "Q1 should match")
		assert.Equal(t, "Savior_DT", history.Q2, "Q2 should match")
		assert.Equal(t, "test", history.Q3, "Q3 should match")
		assert.Equal(t, uint(0), *history.UserID, "UserID should be nil for unauthenticated request")
	})
}
