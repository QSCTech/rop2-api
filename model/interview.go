package model

import "time"

type InterviewStatus int8

const (
	//默认，人满就关
	Auto InterviewStatus = 0
	//无限容量，Capacity无效
	UnlimitedCapacity InterviewStatus = 10
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
	//面试备注（所有人可见，包括未选择此面试的）
	Comment *string `json:"comment" gorm:"type:text;default:null"`

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

//面试基本信息，包括已用容量(JOIN获得结果)。没有createAt和updateAt
//供候选人报名时查看
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
	Comment  *string   `json:"comment" gorm:"type:text;default:null"`

	UsedCapacity int32 `json:"usedCapacity"`
}

//获取指定表单、阶段下，一个或多个部门的面试数组
func GetInterviews(formId uint32, departs []uint32, step StepType) []*InterviewInfo {
	interviews := make([]*InterviewInfo, 0)
	db.
		Select("interviews.*", "COUNT(interview_schedules.zju_id) AS UsedCapacity").
		Model(&Interview{}). //from子句
		Joins("LEFT JOIN interview_schedules ON (interviews.id = interview_schedules.interview)").
		Where("interviews.form = ? AND interviews.depart IN ? AND interviews.step = ?", formId, departs, step).
		Group("interviews.id").
		Scan(&interviews)
	return interviews
}

func AddInterview(formId, depart uint32, step StepType, capacity int32, location string, startAt, endAt time.Time, comment *string) uint32 {
	obj := Interview{
		Form:     formId,
		Depart:   depart,
		Step:     step,
		Capacity: capacity,
		Status:   0, //默认状态
		Location: location,
		StartAt:  startAt,
		EndAt:    endAt,
		Comment:  comment,
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
