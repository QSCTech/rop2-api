package model

import "time"

//候选人的某个志愿
type Intent struct {
	//申请人学号
	ZjuId string `json:"zjuId" gorm:"type:char(10);primaryKey"`
	//来源表单（也确定了申请的组织）
	Form uint32 `json:"form" gorm:"primaryKey"`
	//志愿部门，可能为默认部门（如果未选择志愿部门）
	Depart uint32 `json:"depart" gorm:"primaryKey"`
	//当前所在阶段。1~127=第n阶段(可重命名)
	Step StepType `json:"step" gorm:"not null;default:0"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
}
