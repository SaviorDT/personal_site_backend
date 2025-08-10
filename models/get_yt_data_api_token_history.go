package models

import (
	"gorm.io/gorm"
)

type YTDataAPITokenHistory struct {
	gorm.Model `gorm:"embedded"` // ID, CreatedAt, UpdatedAt
	Q1         string            `gorm:"size:256;not null"`
	Q2         string            `gorm:"size:256;not null"`
	Q3         string            `gorm:"size:256;not null"`
	UserID     *uint             `gorm:"index"` // 可以為 null，關聯到 User
}
