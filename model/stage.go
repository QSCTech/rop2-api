package model

import "time"

type Stage struct {
	Id    uint32  `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name  string  `json:"name" gorm:"type:varchar(40);not null;uniqueIndex:uni_name_owner"` //须在组织内唯一的阶段名
	Tasks string  `json:"tasks" gorm:"not null;type:json"`
	Next  *uint32 `json:"next"` //可空

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	Owner uint32 `json:"owner" gorm:"not null;uniqueIndex:uni_name_owner"`
}
