package model

import (
	"errors"
	"regexp"
	"time"

	"gorm.io/gorm"
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

// 查询向某个表单递交过答卷的所有Person信息
func GetParticipants(formId uint32) *[]Person {
	result := make([]Person, 0)
	db.Select("people.*").Table("results").
		Joins("INNER JOIN people ON results.zju_id = people.zju_id").
		Where("results.form = ?", formId).
		Scan(&result)
	return &result
}

// 创建自然人信息，如果已存在则更新
func EnsurePerson(zjuId string, name string) {
	var count int64
	db.Model(&Person{}).Where("zju_id = ?", zjuId).Count(&count)
	if count == 0 {
		db.Create(&Person{ZjuId: zjuId, Name: name, Phone: nil})
	}
}

func SaveProfile(zjuId string, phone string) error {
	if matched, _ := regexp.MatchString((`^\d{11}$`), phone); !matched {
		return errors.New("手机号格式错误")
	}
	return db.Model(&Person{}).Where("zju_id = ?", zjuId).Update("Phone", phone).Error
}

func FindPerson(zjuId string) *Person {
	var person Person
	if err := db.Where("zju_id = ?", zjuId).First(&person).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return &person
}
