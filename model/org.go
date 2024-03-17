package model

import (
	"time"
)

type Org struct {
	Id   uint32 `json:"id" gorm:"primaryKey;autoIncrement;<-:false"` //主键，自动递增
	Name string `json:"name" gorm:"type:varchar(80);not null;unique"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	DefaultDepart uint32 `json:"defaultDepart" gorm:"uniqueIndex"`
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
