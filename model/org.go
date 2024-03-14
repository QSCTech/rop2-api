package model

import (
	"time"
)

type Org struct {
	Id       uint32    `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `json:"name" gorm:"type:varchar(80);not null;unique"`
	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	DefaultDepart uint32 `json:"defaultDepart" gorm:"uniqueIndex"`
}

type Depart struct {
	Id       uint32    `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `json:"name" gorm:"type:varchar(80);not null"`
	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	Parent uint32 `json:"parent" gorm:"not null"` //归属组织id
}

type User struct {
	ZjuId    string `json:"zjuId" gorm:"type:char(10);primaryKey"`                                //学号，用于唯一识别和登录(复合主键)。考虑到有0开头学号用字符串存
	At       uint32 `json:"at" gorm:"primaryKey;autoIncrement:false;uniqueIndex:uni_nickname_at"` //管理的组织id(复合主键)，防止gorm对int主键添加autoIncrement
	Nickname string `json:"nickname" gorm:"type:char(40);not null;uniqueIndex:uni_nickname_at"`   //在该组织的昵称，在该组织内唯一
	Perm     string `json:"perm" gorm:"type:json;not null"`                                       //权限，json

	CreateAt time.Time `json:"createAt" gorm:"not null"` //创建管理时间
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
