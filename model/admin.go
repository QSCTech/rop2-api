package model

import (
	"time"
)

type Admin struct {
	//学号，用于唯一识别和登录(复合主键)。考虑到有0开头学号用字符串存
	ZjuId string `json:"zjuId" gorm:"type:char(10);primaryKey"`
	//管理的组织id(复合主键)
	At uint32 `json:"at" gorm:"primaryKey;autoIncrement:false;uniqueIndex:uni_nickname_at"`
	//在该组织的昵称，在该组织内唯一
	Nickname string `json:"nickname" gorm:"type:char(40);not null;uniqueIndex:uni_nickname_at"`
	//权限级别
	Level PermLevel `json:"level" gorm:"not null"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
}

// 权限级别，管理员可在同组织下添加权限不高于自己的管理员
type PermLevel int8

const (
	Null       PermLevel = 10 * (iota) //无权限(0)，保留备用
	Observer                           //可查看简历、表单、面试
	Maintainer                         //所有权限。可修改表单、面试、审批
	//为简化，我们减少了权限级别数，估计可查看/可管理两级足以满足大部分需求
)

func GetAdmin(zjuId string, at uint32) *Admin {
	var pobj = &Admin{}
	result := db.First(pobj, "zju_id = ? AND at = ?", zjuId, at)
	if result.Error != nil {
		return nil
	} else {
		return pobj
	}
}
