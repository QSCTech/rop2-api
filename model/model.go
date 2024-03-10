package model

import "time"

type Org struct {
	Id       uint32    `gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `gorm:"not null;type:char(100)"` //一般不作为搜索条件，不设成index
	CreateAt time.Time `gorm:"autoCreateTime"`
}
