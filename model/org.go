package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Org struct {
	Id       uint32    `json:"id" gorm:"primaryKey;autoIncrement;<-:false"` //主键，自动递增
	Name     string    `json:"name" gorm:"type:varchar(80);not null;unique"`
	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	DefaultDepart uint32 `json:"defaultDepart" gorm:"uniqueIndex"`
}

type Depart struct {
	Id       uint32    `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `json:"name" gorm:"type:varchar(80);not null;uniqueIndex:uni_name_parent"` //部门名称，须在组织内唯一
	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	Parent uint32 `json:"parent" gorm:"not null;uniqueIndex:uni_name_parent"` //归属组织id
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

func GetOrgDeparts(orgId uint32) []Depart {
	result := make([]Depart, 0)
	db.Select("Id", "Name", "CreateAt").Find(&result, "parent = ?", orgId)
	return result
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

func CreateDepart(orgId uint32, name string) bool {
	d := &Depart{
		Name:   name,
		Parent: orgId,
	}
	result := db.Select("Name", "Parent").Create(d)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return false
		} else {
			panic(result.Error)
		}
	}
	return true
}

func DeleteDepart(id uint32) {
	db.Delete(&Depart{}, id)
}

func RenameDepart(id uint32, newName string) bool {
	result := db.Model(&Depart{}).Where(id).Update("name", newName)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return false
		} else {
			panic(result.Error)
		}
	}
	return true
}
