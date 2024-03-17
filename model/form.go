package model

import "time"

type Form struct {
	Id       uint32 `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Name     string `json:"name" gorm:"type:varchar(100);not null;uniqueIndex:uni_name_owner"` //须在组织内唯一的表单名称
	Desc     string `json:"desc"`
	Entry    uint32 `json:"entry" gorm:"not null"`
	Children string `json:"children" gorm:"not null;type:json"`
	Enter    uint32 `json:"enter" gorm:"not null"`

	StartAt *time.Time `json:"startAt"` //可空
	EndAt   *time.Time `json:"endAt"`   //可空

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`

	Owner uint32 `json:"owner" gorm:"not null;uniqueIndex:uni_name_owner"`
}

func GetForms(owner uint32) []*Form {
	result := make([]*Form, 0)
	db.
		Order("id desc").
		Select("Id", "Name", "StartAt", "EndAt", "CreateAt", "UpdateAt").
		Find(&result, "owner = ?", owner)
	return result
}

func GetFormDetail(owner uint32, id uint32) *Form {
	pobj := &Form{}
	result := db.First(pobj, "id = ? AND owner = ?", id, owner)
	if result.Error != nil {
		return nil
	}
	return pobj
}

func SaveForm(obj *Form) {
	db.Save(obj)
}
