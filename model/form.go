package model

import (
	"errors"
	"rop2-api/utils"
	"time"
)

type Form struct {
	Id   uint32 `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name string `json:"name" gorm:"type:varchar(100);not null;uniqueIndex:uni_name_owner"` //须在组织内唯一的表单名称
	Desc string `json:"desc"`

	//入口题目组约定为id:1的题目组，不再保留Entry属性

	Children string `json:"children" gorm:"not null;type:json"`

	StartAt *time.Time `json:"startAt"` //可空
	EndAt   *time.Time `json:"endAt"`   //可空

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`

	Owner uint32 `json:"owner" gorm:"not null;uniqueIndex:uni_name_owner"`
}

// 按id降序，查询指定组织所有表单的简略信息
func GetForms(owner uint32) []*Form {
	result := make([]*Form, 0)
	db.
		Order("id desc").
		Select("Id", "Name", "StartAt", "EndAt", "CreateAt", "UpdateAt").
		Find(&result, "owner = ?", owner)
	return result
}

// 检查指定表单id是否为指定组织所创建
func CheckOwner(owner uint32, id uint32) bool {
	var count int64
	db.Table("forms").Where("id = ? AND owner = ?", id, owner).Count(&count)
	return count > 0
}

// 查询表单详情同时限定owner，适用于管理员查询
func GetFormDetail(owner uint32, id uint32) *Form {
	pobj := &Form{}
	result := db.First(pobj, "id = ? AND owner = ?", id, owner)
	if result.Error != nil {
		return nil
	}
	return pobj
}

// 根据id查询表单详情，仅部分字段(不包含CreateAt,UpdateAt)
func ApplicantGetFormDetail(id uint32) *Form {
	pobj := &Form{}
	result := db.Select("Id", "Name", "Desc", "Children", "StartAt", "EndAt", "Owner").First(pobj, "id = ?", id)
	if result.Error != nil {
		return nil
	}
	return pobj
}

type FormUpdate struct {
	Id       uint32  `json:"id"`
	Name     *string `json:"name"` //须在组织内唯一的表单名称
	Desc     *string `json:"desc"`
	Children *string `json:"children"`

	StartAt *time.Time `json:"startAt"`
	EndAt   *time.Time `json:"endAt"`
}

// 修改指定表单
func SaveForm(obj FormUpdate) error {
	updateMap := make(map[string]interface{})
	if obj.Name != nil {
		if diff := utils.LenBetween(*obj.Name, 1, 25); diff != 0 {
			if diff > 0 {
				return errors.New("标题过长")
			} else {
				return errors.New("标题过短")
			}
		}
		updateMap["name"] = obj.Name
	}
	if obj.Desc != nil {
		if diff := utils.LenBetween(*obj.Desc, 0, 200); diff != 0 {
			return errors.New("简介过长")
		}
		updateMap["Desc"] = obj.Desc
	}
	if obj.Children != nil {
		updateMap["Children"] = obj.Children
	}
	if obj.StartAt != nil {
		//为unix时间戳<100即为设空，为nil保持不变
		if obj.StartAt.Before(time.Unix(100, 0)) {
			updateMap["Start_At"] = nil
		} else {
			updateMap["Start_At"] = obj.StartAt
		}
	}
	if obj.EndAt != nil {
		//>2048年即为设空，为nil保持不变
		if obj.EndAt.After(time.Date(2048, 1, 1, 0, 0, 0, 0, time.Local)) {
			updateMap["End_At"] = nil
		} else {
			updateMap["End_At"] = obj.EndAt
		}
	}
	db.Table("forms").Where("id = ?", obj.Id).Updates(updateMap)
	return nil
}

func CreateForm(owner uint32, name string) (uint32, error) {
	form := &Form{
		Name:     name,
		Owner:    owner,
		Children: `[{"id":1,"children":[],"label":"基本信息"}]`, //TODO 修改新问卷的默认问题
	}
	result := db.Select("Name", "Owner", "Children").Create(form)
	if result.Error != nil {
		return 0, result.Error
	}
	return form.Id, nil
}

// 删除表单，限定owner，受影响行数>0返回true
func DeleteForm(owner uint32, formId uint32) bool {
	result := db.Delete(&Form{}, "id = ? AND owner = ?", formId, owner)
	return result.RowsAffected > 0
}
