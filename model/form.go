package model

import "time"

type Form struct {
	Id       uint32 `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string `json:"name" gorm:"not null"`
	Desc     string `json:"desc" gorm:"not null"`
	Entry    uint32 `json:"entry" gorm:"not null"`
	Children string `json:"children" gorm:"not null;type:json"`

	StartAt time.Time
	EndAt   time.Time

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`

	Owner uint32 `json:"owner" gorm:"not null"`
}
