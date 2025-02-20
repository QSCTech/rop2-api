package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TODO: 把类型为string的zjuId改为PersonId方便后面修改
type PersonId = string

// 自然人信息，在所有表单都共用。其中，学号为唯一标识
// 统一认证登录后可以自行修改信息（所有表单都将同步修改）。
// 目前管理员不适用这些信息。
type Person struct {
	ZjuId PersonId `json:"zjuId" gorm:"type:char(10);primaryKey"`          //学号
	Name  string   `json:"name" gorm:"type:char(20);not null"`             //真实姓名，不可空
	Phone *string  `json:"phone" gorm:"type:char(11);unique;default:null"` //手机号，唯一
	//暂时不添加邮箱、QQ、微信等字段。表单可以自行添加这些题目

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

// 创建自然人信息，如果已存在则更新姓名(避免批量导入时存在错误)
func EnsurePerson(zjuId string, name string) {
	db.
		Select("zju_id", "name").
		Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{"name"}),
		}).
		Create(&Person{ZjuId: zjuId, Name: name})
}

func FindPerson(zjuId string) *Person {
	var person Person
	if err := db.Where("zju_id = ?", zjuId).First(&person).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return &person
}
