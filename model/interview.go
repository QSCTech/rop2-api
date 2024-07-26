package model

import "time"

type InviewStatus int8

const (
	//默认，人满就关
	Auto InviewStatus = 0
	//无限容量，Capacity无效
	Unlimited InviewStatus = 10
	//被管理员手动冻结，不可报名/不可取消
	Frozen InviewStatus = 20
)

type Interview struct {
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"`

	Form     uint32       `json:"form" gorm:"not null"`
	Depart   uint32       `json:"depart" gorm:"not null"`
	Step     StepType     `json:"step" gorm:"not null"`
	Capacity int32        `json:"capacity" gorm:"not null"`
	Status   InviewStatus `json:"status" gorm:"not null;default:0"`

	//面试详情
	Location string    `json:"location" gorm:"not null"`
	StartAt  time.Time `json:"startAt" gorm:"not null"`
	EndAt    time.Time `json:"endAt" gorm:"not null"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

func GetInterview(id uint32) *Interview {
	var obj Interview
	if db.First(&obj, id).Error != nil {
		return nil
	} else {
		return &obj
	}
}

func GetInterviews(formId uint32, departs []uint32, step StepType) []Interview {
	var interviews []Interview
	db.Where("form = ? AND depart IN ? AND step = ?", formId, departs, step).Find(&interviews)
	return interviews
}

func AddInterview(formId, depart uint32, step StepType, capacity int32, location string, startAt, endAt time.Time) uint32 {
	obj := Interview{
		Form:     formId,
		Depart:   depart,
		Step:     step,
		Capacity: capacity,
		Status:   0, //默认状态
		Location: location,
		StartAt:  startAt,
		EndAt:    endAt,
	}
	db.Create(&obj)
	return obj.Id
}

func DeleteInterview(id uint32) error {
	return db.Delete(&Interview{}, id).Error
}
