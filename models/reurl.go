package models

import (
    "time"

    "gorm.io/gorm"
)

// Reurl represents a redirect mapping created by a user.
// Key: unique key used to look up the mapping (can contain unicode characters)
// TargetURL: destination URL to redirect to
// ExpiresAt: optional expiration time; nil means never expire
type Reurl struct {
    gorm.Model `gorm:"embedded"`
    Key       string     `gorm:"size:256;not null;uniqueIndex" json:"key"`
    TargetURL string     `gorm:"size:2048;not null" json:"target_url"`
    ExpiresAt *time.Time `json:"expires_at"`
    OwnerID   uint       `gorm:"not null;index" json:"owner_id"`
    Owner     User       `gorm:"foreignKey:OwnerID" json:"owner"`
}

func (Reurl) TableName() string {
    return "reurls"
}
