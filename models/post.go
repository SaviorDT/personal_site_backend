package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// PostStatus 定義文章狀態
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"     // 草稿
	PostStatusPublished PostStatus = "published" // 已發布
	PostStatusArchived  PostStatus = "archived"  // 封存
)

// PostVisibility 定義文章可見性
type PostVisibility string

const (
	PostVisibilityPublic  PostVisibility = "public"  // 公開
	PostVisibilityPrivate PostVisibility = "private" // 私人
)

// Post 文章模型
type Post struct {
	gorm.Model
	AuthorID    uint           `gorm:"not null;index"`                                // 作者 ID
	Author      User           `gorm:"foreignKey:AuthorID"`                           // 作者關聯
	Title       string         `gorm:"size:255;not null"`                             // 標題
	Slug        string         `gorm:"size:255;uniqueIndex"`                          // URL 友善名稱（用於前端靜態文章）
	Content     string         `gorm:"type:text;not null"`                            // 內容
	Summary     string         `gorm:"size:500"`                                      // 摘要
	CoverImage  string         `gorm:"size:500"`                                      // 封面圖
	Status      PostStatus     `gorm:"size:20;not null;default:'draft'"`              // 狀態
	Visibility  PostVisibility `gorm:"size:20;not null;default:'public'"`             // 可見性
	ViewCount   int            `gorm:"default:0"`                                     // 瀏覽次數
	PublishedAt *time.Time     `gorm:"index"`                                         // 發布時間
	Tags        []Tag          `gorm:"many2many:post_tags;"`                          // 標籤（多對多）
	Comments    []Comment      `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"` // 留言
	Reactions   []Reaction     `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE"` // 反應
}

// Tag 標籤模型
type Tag struct {
	gorm.Model
	Name  string `gorm:"size:50;not null;uniqueIndex"` // 標籤名稱
	Slug  string `gorm:"size:50;not null;uniqueIndex"` // URL 友善名稱
	Color string `gorm:"size:7;default:'#3B82F6'"`     // 標籤顏色
	Posts []Post `gorm:"many2many:post_tags;"`         // 文章（多對多）
}

func (Post) TableName() string { return "posts" }

// BeforeCreate 在創建前生成 Slug
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		// 生成基礎 slug (保留中文，只替換空格)
		slug := strings.TrimSpace(p.Title)
		slug = strings.ReplaceAll(slug, " ", "-")

		// 加上短隨機字串以確保唯一性
		rand.Seed(time.Now().UnixNano())
		randomSuffix := fmt.Sprintf("%06d", rand.Intn(1000000))
		p.Slug = fmt.Sprintf("%s-%s", slug, randomSuffix)
	}
	return nil
}

// BeforeSave 驗證資料
func (p *Post) BeforeSave(tx *gorm.DB) error {
	// 如果狀態變更為已發布且 PublishedAt 為空，設定發布時間
	if p.Status == PostStatusPublished && p.PublishedAt == nil {
		now := time.Now()
		p.PublishedAt = &now
	}
	return nil
}

// FindPostByIDOrSlug 根據 ID 或 Slug 查找文章
// 如果 identifier 可以解析為數字，則用 ID 查詢；否則用 Slug 查詢
func FindPostByIDOrSlug(db *gorm.DB, identifier string) (*Post, error) {
	var post Post

	// 嘗試解析為數字
	if _, err := strconv.Atoi(identifier); err == nil {
		// 如果是數字，嘗試用 ID 查詢
		if err := db.First(&post, identifier).Error; err == nil {
			return &post, nil
		}
	}

	// 如果 ID 查詢失敗或不是數字，用 Slug 查詢
	if err := db.Where("slug = ?", identifier).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
