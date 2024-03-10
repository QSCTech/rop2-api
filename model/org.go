package model

import "time"

type Org struct {
	Id       uint32    `gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `gorm:"not null;type:char(100)"`
	CreateAt time.Time `gorm:"not null;autoCreateTime"`

	DefaultDepart uint32 `gorm:"not null;default:0"` //default:0保证不存在对应的depart 插入后需立即设置为有效的默认部门
}

type Depart struct {
	Id       uint32    `gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `gorm:"not null;type:char(100)"`
	CreateAt time.Time `gorm:"not null;autoCreateTime"`
	OrgId    uint32    `gorm:"not null"` //外键命名约定
	Org      Org
}
