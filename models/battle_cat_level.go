package models

type BattleCatLevel struct {
	ID      uint   `gorm:"primaryKey"`
	Stage   string `gorm:"type:char(3);not null;Index"`
	Level   string `gorm:"size:16;not null;Index"`
	Name    string `gorm:"size:128;not null;Index"`
	HP      uint16 `gorm:"not null"`
	Enemies string `gorm:"size:64;not null;Index"`
}
