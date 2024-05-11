package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Depart struct {
	Id       uint32    `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `json:"name" gorm:"type:varchar(80);not null;uniqueIndex:uni_name_owner"` //部门名称，须在组织内唯一
	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	Owner uint32 `json:"owner" gorm:"not null;uniqueIndex:uni_name_owner"` //归属组织id
}

func GetOrgDeparts(orgId uint32) *[]Depart {
	result := make([]Depart, 0)
	db.Select("Id", "Name", "CreateAt").Find(&result, "owner = ?", orgId)
	return &result
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

func CreateDepart(orgId uint32, name string) (bool, *Depart) {
	d := &Depart{
		Name:  name,
		Owner: orgId,
	}
	result := db.Select("Name", "owner").Create(d)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return false, nil
		} else {
			panic(result.Error)
		}
	}
	return true, d
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
