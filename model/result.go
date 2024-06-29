package model

import "time"

//报名者答卷。同一人报名多个志愿部门，仅一份答卷，但生成多个Intent。
type Result struct {
	Form  uint32 `json:"from" gorm:"primaryKey"`
	ZjuId string `json:"zjuId" gorm:"type:char(10);primaryKey"`
	//实际答卷内容
	Content string `json:"content" gorm:"not null;type:json"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

//创建/更新答卷
func SaveResult(formId uint32, zjuId string, content string) error {
	result := &Result{
		Form:    formId,
		ZjuId:   zjuId,
		Content: content,
	}
	return db.Save(result).Error
}
