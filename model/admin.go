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

func GetAdmin(zjuId string, at *uint32) []*Admin {
	var pobj = make([]*Admin, 0)
	if at == nil {
		db.Find(&pobj, "zju_id = ?", zjuId)
	} else {
		db.Find(&pobj, "zju_id = ? and at = ?", zjuId, *at)
	}
	return pobj
}

type AdminChoice struct {
	OrgId   uint32 `json:"orgId"`
	OrgName string `json:"orgName"`
}

func GetAvailableOrgs(zjuId string) []*AdminChoice {
	profiles := make([]*AdminChoice, 2)
	db.Table("admins").Select("orgs.name as OrgName", "admins.at as OrgId").Joins("JOIN orgs ON admins.at = orgs.id").Where("admins.zju_id = ?", zjuId).Scan(&profiles)
	return profiles
}

type AdminProfile struct {
	Nickname  string    `json:"nickname"`
	Zju_Id    string    `json:"zjuId"`
	Level     PermLevel `json:"level"`
	Create_At time.Time `json:"createAt"`
}

type AdminList struct {
	Admins        []*AdminProfile `json:"admins"`
	Count         int64           `json:"count"`
	FilteredCount int64           `json:"filteredCount"`
}

func GetAdminsInOrg(orgId uint32, offset int, limit int, filter string) AdminList {
	var count int64 //指定组织下管理员总数
	db.Table("admins").Where("at = ?", orgId).Count(&count)

	if filter == "" {
		filter = "^" //匹配所有
	}
	admins := make([]*AdminProfile, 0)
	db.
		Table("admins").
		Select("Nickname", "Level", "Create_At", "Zju_Id").
		Where("at = ?", orgId).
		Where("nickname REGEXP ?", filter).
		Order("Create_At DESC"). //按创建时间降序
		Offset(offset).
		Limit(limit).
		Scan(&admins)

	var filteredCount int64
	db.
		Table("admins").
		Where("at = ?", orgId).
		Where("nickname REGEXP ?", filter).
		Count(&filteredCount) //mysql有优化

	var list = AdminList{Admins: admins, Count: count, FilteredCount: filteredCount}
	return list
}

func SetAdmin(at uint32, zjuId string, nickname string, level PermLevel) {
	if level <= Null {
		db.Delete(&Admin{}, "at = ? and zju_id = ?", at, zjuId)
		return
	}

	var admin = &Admin{At: at, ZjuId: zjuId, Nickname: nickname, Level: level}
	db.Save(admin)
}
