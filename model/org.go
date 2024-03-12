package model

import (
	"time"
)

type Org struct {
	Id       uint32    `gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `gorm:"not null"`
	CreateAt time.Time `gorm:"not null;autoCreateTime"`

	DefaultDepart uint32 `gorm:"uniqueIndex"`
}

type Depart struct {
	Id       uint32    `gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `gorm:"not null"`
	CreateAt time.Time `gorm:"not null;autoCreateTime"`

	Parent uint32 `gorm:"not null"`
}

func GetOrg(id uint32) *Org {
	var pobj = &Org{}
	result := db.First(pobj, id)
	if result.Error != nil {
		return nil
	} else {
		return pobj
	}
}

func GetDepart(id uint32) *Depart {
	var pobj = &Depart{}
	result := db.First(pobj, id)
	if result.Error != nil {
		return nil
	} else {
		return pobj
	}
}
