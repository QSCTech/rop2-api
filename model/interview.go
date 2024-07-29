package model

import "time"

type InterviewStatus int8

const (
	//默认，人满就关
	Auto InterviewStatus = 0
	//无限容量，Capacity无效
	Unlimited InterviewStatus = 10
	//被管理员手动冻结，不可报名/不可取消
	Frozen InterviewStatus = 20
)

type Interview struct {
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"`

	Form     uint32          `json:"form" gorm:"not null"`
	Depart   uint32          `json:"depart" gorm:"not null"`
	Step     StepType        `json:"step" gorm:"not null"`
	Capacity int32           `json:"capacity" gorm:"not null"`
	Status   InterviewStatus `json:"status" gorm:"not null;default:0"`

	//面试详情
	Location string    `json:"location" gorm:"not null"`
	StartAt  time.Time `json:"startAt" gorm:"not null"`
	EndAt    time.Time `json:"endAt" gorm:"not null"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

func GetInterviewById(id uint32) *Interview {
	var obj Interview
	if db.First(&obj, id).Error != nil {
		return nil
	} else {
		return &obj
	}
}

//面试基本信息，包括已用容量。没有createAt和updateAt
type InterviewInfo struct {
	Id uint32 `json:"id"`

	Depart   uint32          `json:"depart"`
	Step     StepType        `json:"step"`
	Capacity int32           `json:"capacity"`
	Status   InterviewStatus `json:"status"`

	//面试详情
	Location string    `json:"location"`
	StartAt  time.Time `json:"startAt"`
	EndAt    time.Time `json:"endAt"`

	UsedCapacity int32 `json:"usedCapacity"`
}

func GetInterviews(formId uint32, departs []uint32, step StepType) []InterviewInfo {
	var interviews []InterviewInfo = make([]InterviewInfo, 0)
	db.
		Select(
			"interviews.id", "interviews.depart", "interviews.step", "interviews.capacity", "interviews.status",
			"interviews.location", "interviews.start_at", "interviews.end_at",
					"COUNT(interview_schedules.zju_id) AS UsedCapacity").
		Table("interviews"). //from子句
		Joins("LEFT JOIN interview_schedules ON (interviews.id = interview_schedules.interview)").
		Where("interviews.form = ? AND interviews.depart IN ? AND interviews.step = ?", formId, departs, step).
		Group("interviews.id").
		Scan(&interviews)
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

func FreezeInterview(id uint32) error {
	return db.Model(&Interview{}).Where("id = ?", id).Update("status", Frozen).Error
}
