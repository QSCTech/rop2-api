package model

import "time"

type Form struct {
	Id       uint32 `gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string `gorm:"not null"`
	Desc     string `gorm:"not null"`
	Entry    uint32 `gorm:"not null"`
	Children string `gorm:"not null;type:json"`

	StartAt time.Time
	EndAt   time.Time

	CreateAt time.Time `gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `gorm:"not null;autoUpdateTime"`

	Owner uint32 `gorm:"not null"`
}
