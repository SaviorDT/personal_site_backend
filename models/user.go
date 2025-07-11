package models

import (
	"fmt"
	"time"
	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleUser, RoleGuest:
		return true
	default:
		return false
	}
}

type AuthProvider string

const (
	AuthProviderPassword AuthProvider = "password"
	AuthProviderGitHub   AuthProvider = "github"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderLine     AuthProvider = "line"
)

func (a AuthProvider) IsValid() bool {
	switch a {
	case AuthProviderPassword, AuthProviderGitHub, AuthProviderGoogle, AuthProviderLine:
		return true
	default:
		return false
	}
}

type User struct {
	Model      Model        `gorm:"embedded;"` // ID, CreatedAt, UpdatedAt
	Nickname   string       `gorm:"size:64;not null"`
	Role       Role         `gorm:"size:32;not null"`
	Provider   AuthProvider `gorm:"size:16;not null"`
	Email      string       `gorm:"size:128;not null;unique"`
	Identifier string       `gorm:"size:256"` // hashed password, or provider id
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	if !u.Provider.IsValid() {
		return fmt.Errorf("invalid auth provider: %s", u.Provider)
	}
	if !u.Role.IsValid() {
		return fmt.Errorf("invalid role: %s", u.Role)
	}
	return nil
}
