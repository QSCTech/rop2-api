package model

import "time"

type Candidate struct {
	Id    uint32 `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	ZjuId string `json:"zjuId" gorm:"not null;type:char(10)"`
	Of    uint32 `json:"of" gorm:"not null"`    //志愿部门信息（每个部门独立考核）
	Form  uint32 `json:"form" gorm:"not null"`  //来源表单
	Stage uint32 `json:"stage" gorm:"not null"` //此条数据所在阶段

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	Gone *uint32 `json:"gone"` //候选人已进入的下一候选数据id，可空
}
