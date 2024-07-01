package model

import (
	"time"

	"gorm.io/gorm"
)

// 候选人的某个志愿
type Intent struct {
	//自增主键，志愿id
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"`
	//来源表单（也确定了申请的组织）
	Form uint32 `json:"form" gorm:"not null;uniqueIndex:zjuid_form_intent,sort:desc"`
	//申请人学号
	ZjuId string `json:"zjuId" gorm:"type:char(10);not null;uniqueIndex:zjuid_form_intent,sort:desc"`
	//志愿部门，可能为默认部门（如果未选择志愿部门）
	Depart uint32 `json:"depart" gorm:"not null;uniqueIndex:zjuid_form_intent"`
	//志愿排序。1~127=第n志愿
	Order int8 `json:"order" gorm:"not null;uniqueIndex:zjuid_form_intent,sort:asc"`
	//当前所在阶段。1~127=第n阶段(可重命名)
	Step StepType `json:"step" gorm:"not null;default:0"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

func SaveIntents(formId uint32, zjuId string, intentDeparts []uint32) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Intent{}, "form = ? AND zju_id = ?", formId, zjuId).Error; err != nil {
			return err
		}
		intents := make([]Intent, len(intentDeparts))
		for i, v := range intentDeparts {
			intents[i] = Intent{
				Form:   formId,
				ZjuId:  zjuId,
				Depart: v,
				Order:  int8(i + 1),
			}
		}
		return tx.Select("Form", "ZjuId", "Depart", "Order").Create(intents).Error
	})
}
