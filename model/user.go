package model

import (
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ZjuId    string `json:"zjuId" gorm:"type:char(10);primaryKey"`                                //学号，用于唯一识别和登录(复合主键)。考虑到有0开头学号用字符串存
	At       uint32 `json:"at" gorm:"primaryKey;autoIncrement:false;uniqueIndex:uni_nickname_at"` //管理的组织id(复合主键)，防止gorm对int主键添加autoIncrement
	Nickname string `json:"nickname" gorm:"type:char(40);not null;uniqueIndex:uni_nickname_at"`   //在该组织的昵称，在该组织内唯一
	Perm     string `json:"perm" gorm:"type:json;not null"`                                       //权限，json {"部门id":级别,...}

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"` //创建管理时间
}

// 权限级别，管理员可在同一部门下添加权限不高于自己的管理员
type PermLevel uint32

const (
	Null       PermLevel = 100 * (iota) //无权限(0)，保留备用
	Observer                            //可查看简历、表单、面试
	Insepector                          //可修改表单、面试、审批
	Maintainer                          //可修改部门、阶段
)

// 检查权限级别是否有效(授权时使用)
func IsValidLevel(level PermLevel, granterLevel PermLevel) bool {
	return level >= Observer && level <= Maintainer && level%100 == 0 && level <= granterLevel
}

// 部门-权限的映射表
type PermMap map[uint32]PermLevel

// 从json中解析部门-权限映射
func ParsePerm(perm string) (result PermMap) {
	result = make(PermMap)
	iter := jsoniter.ParseString(jsoniter.ConfigFastest, perm)
	for depIdStr := iter.ReadObject(); depIdStr != ""; depIdStr = iter.ReadObject() {
		depId, err := strconv.ParseUint(depIdStr, 10, 32)
		if err != nil {
			panic(err)
		}
		result[uint32(depId)] = PermLevel(iter.ReadUint32())
	}
	return result
}

// 获取在指定部门的权限，考虑默认部门
func GetLevel(permMap PermMap, departId uint32, defaultDepartId uint32) PermLevel {
	defaultDepPerm, depPerm := permMap[defaultDepartId], permMap[departId] //不存在默认为0
	return max(defaultDepPerm, depPerm)
}

func GetUser(zjuId string, at uint32) *User {
	var pobj = &User{}
	result := db.First(pobj, "zju_id = ? AND at = ?", zjuId, at)
	if result.Error != nil {
		return nil
	} else {
		return pobj
	}
}

// 检查用户是否在默认部门有足够的权限
func HasOrgLevel(perm string, orgId uint32, requireLevel PermLevel) bool {
	org := GetOrg(orgId)
	permMap := ParsePerm(perm)
	orgPerm := permMap[org.DefaultDepart]
	return orgPerm >= requireLevel
}
