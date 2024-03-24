package model

import "time"

type StepType int8

//将阶段进行重命名或配置其通知
type Stage struct {
	//生效部门，可以为默认部门（为所有其它部门提供默认名称和通知设定）
	Owner uint32 `json:"owner" gorm:"primaryKey;uniqueIndex:uni_owner_name"`
	//阶段序号。从1开始(表示“第一阶段”)
	Step StepType `json:"step" gorm:"primaryKey"`
	//新名字，需要在命名的部门内唯一。可以为空，表示不进行而使用默认递增命名
	Name *string `json:"name" gorm:"type:varchar(40);uniqueIndex:uni_owner_name"`

	//进入时发送的通知，可以为空，表示不发送。
	OnEnter *uint32 `json:"onEnter"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

const (
	Invalid  StepType = -2  //已失效，不发送拒信的已拒绝
	Rejected StepType = -1  //已拒绝
	Accepted StepType = -50 //已录取
	Applied  StepType = 0   //已填表，下一阶段为“第一阶段”
)
