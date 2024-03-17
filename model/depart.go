package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Depart struct {
	Id       uint32    `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string    `json:"name" gorm:"type:varchar(80);not null;uniqueIndex:uni_name_parent"` //部门名称，须在组织内唯一
	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	Parent uint32 `json:"parent" gorm:"not null;uniqueIndex:uni_name_parent"` //归属组织id
}

func GetOrgDeparts(orgId uint32) []*Depart {
	result := make([]*Depart, 0)
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

func DeleteDepart(id uint32) bool {
	result := db.Delete(&Depart{}, id)
	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			//默认部门受外键约束不能删除
			return false
		}
		panic(err)
	}
	return true
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
