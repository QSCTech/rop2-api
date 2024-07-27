package model

import "time"

type InterviewSchedule struct {
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"` //面试安排id

	ZjuId     string `json:"zjuId" gorm:"type:char(10)"` //学号
	Interview uint32 `json:"interview" gorm:"not null"`  //面试id

	//用于统计的字段
	Step   StepType `json:"step" gorm:"not null"`   //阶段
	Form   uint32   `json:"form" gorm:"not null"`   //表单id
	Depart uint32   `json:"depart" gorm:"not null"` //部门id

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}
