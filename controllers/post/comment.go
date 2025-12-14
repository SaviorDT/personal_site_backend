package post

import (
	"fmt"
	"net/http"
	"personal_site/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required,min=1"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

func CreateComment(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var post models.Post
	postPtr, err := models.FindPostByIDOrSlug(db, postID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
		return
	}
	post = *postPtr
	if req.ParentID != nil {
		var parent models.Comment
		if err := db.First(&parent, *req.ParentID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent comment not found"})
			return
		}
		if parent.PostID != post.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent comment does not belong to this post"})
			return
		}
	}
	comment := models.Comment{PostID: post.ID, AuthorID: userID.(uint), Content: req.Content, ParentID: req.ParentID}
	if err := db.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}
	db.Preload("Author").First(&comment, comment.ID)
	c.JSON(http.StatusCreated, gin.H{"message": "Comment created successfully", "comment": comment})
}

func GetComments(c *gin.Context, db *gorm.DB) {
	postID := c.Param("id")

	post, err := models.FindPostByIDOrSlug(db, postID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var comments []models.Comment
	if err := db.Where("post_id = ? AND is_deleted = ?", post.ID, false).
		Preload("Author").
		Preload("Reactions").
		Order("created_at ASC"). // Order by time ascending usually makes more sense for flat conversation flow, but DESC is fine if tree sorts it. Keeping DESC for now or matching expectation. Usually comments are old->new. Let's stick to DESC if that was default, or switch to ASC? Default was DESC. Let's keep DESC.
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	fmt.Printf("[Debug] GetComments: fetched %d comments. First one: Content='%s', Author='%+v' (Nickname: %s)\n", len(comments), comments[0].Content, comments[0].Author, comments[0].Author.Nickname)
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func UpdateComment(c *gin.Context, db *gorm.DB) {
	commentID := c.Param("comment_id")
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}
	if comment.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own comments"})
		return
	}
	comment.Content = req.Content
	comment.IsEdited = true
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}
	db.Preload("Author").First(&comment, comment.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Comment updated", "comment": comment})
}

func DeleteComment(c *gin.Context, db *gorm.DB) {
	commentID := c.Param("comment_id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var comment models.Comment
	if err := db.First(&comment, commentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}
	if comment.AuthorID != userID.(uint) {
		// allow post author to delete
		var post models.Post
		if err := db.First(&post, comment.PostID).Error; err == nil {
			if post.AuthorID != userID.(uint) {
				c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this comment"})
				return
			}
		}
	}
	comment.IsDeleted = true
	comment.Content = "[此留言已刪除]"
	if err := db.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}
