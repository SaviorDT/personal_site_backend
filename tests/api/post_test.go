package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"personal_site/models"

	"github.com/stretchr/testify/assert"
)

func TestCreatePost(t *testing.T) {
	router, db := setupTestRouter()
	defer cleanupTestDB(db)

	// 建立測試使用者並登入
	user, token := createTestUserAndLogin(t, db)

	// 測試建立文章
	postData := map[string]interface{}{
		"title":      "測試文章",
		"content":    "這是測試內容",
		"summary":    "測試摘要",
		"status":     "published",
		"visibility": "public",
		"tags":       []string{"測試", "Go"},
	}

	body, _ := json.Marshal(postData)
	req, _ := http.NewRequest("POST", "/api/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Post created successfully", response["message"])
	assert.NotNil(t, response["post"])

	post := response["post"].(map[string]interface{})
	assert.Equal(t, "測試文章", post["title"])
	assert.Equal(t, float64(user.ID), post["author_id"])
}

func TestGetPosts(t *testing.T) {
	router, db := setupTestRouter()
	defer cleanupTestDB(db)

	// 建立測試使用者和文章
	user, _ := createTestUserAndLogin(t, db)
	createTestPost(t, db, user.ID)

	// 測試取得文章列表
	req, _ := http.NewRequest("GET", "/api/posts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	posts := response["posts"].([]interface{})
	assert.Greater(t, len(posts), 0)
}

func TestGetPost(t *testing.T) {
	router, db := setupTestRouter()
	defer cleanupTestDB(db)

	user, _ := createTestUserAndLogin(t, db)
	post := createTestPost(t, db, user.ID)

	// 測試取得單一文章
	req, _ := http.NewRequest("GET", "/api/posts/"+string(rune(post.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	postData := response["post"].(map[string]interface{})
	assert.Equal(t, "測試文章", postData["title"])
}

func TestUpdatePost(t *testing.T) {
	router, db := setupTestRouter()
	defer cleanupTestDB(db)

	user, token := createTestUserAndLogin(t, db)
	post := createTestPost(t, db, user.ID)

	// 測試更新文章
	updateData := map[string]interface{}{
		"title":   "更新後的標題",
		"content": "更新後的內容",
	}

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/api/posts/"+string(rune(post.ID)), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Post updated successfully", response["message"])
}

func TestDeletePost(t *testing.T) {
	router, db := setupTestRouter()
	defer cleanupTestDB(db)

	user, token := createTestUserAndLogin(t, db)
	post := createTestPost(t, db, user.ID)

	// 測試刪除文章
	req, _ := http.NewRequest("DELETE", "/api/posts/"+string(rune(post.ID)), nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: token})

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// 驗證文章已被刪除
	var deletedPost models.Post
	err := db.First(&deletedPost, post.ID).Error
	assert.Error(t, err) // 應該找不到
}

// Helper functions
func createTestPost(t *testing.T, db *gorm.DB, authorID uint) *models.Post {
	post := &models.Post{
		AuthorID:   authorID,
		Title:      "測試文章",
		Content:    "測試內容",
		Status:     models.PostStatusPublished,
		Visibility: models.PostVisibilityPublic,
	}
	err := db.Create(post).Error
	assert.NoError(t, err)
	return post
}
