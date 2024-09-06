package model

import (
	"errors"
	"rop2-api/utils"
	"time"

	"gorm.io/gorm"
)

type InterviewSchedule struct {
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"` //面试安排id

	ZjuId     PersonId `json:"zjuId" gorm:"type:char(10);uniqueIndex:uni_zjuId_interview"` //学号
	Interview uint32   `json:"interview" gorm:"not null;uniqueIndex:uni_zjuId_interview"`  //面试id，同样指定了此次面试安排归属组织、表单、部门、阶段

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

// 查询报名指定面试的所有学号
func GetScheduledIds(interviewId uint32) []PersonId {
	result := make([]PersonId, 0)
	db.Model(&InterviewSchedule{}).Where("interview = ?", interviewId).Pluck("zju_id", &result)
	return result
}

func RemoveScheduledId(interviewId uint32, zjuId PersonId) {
	db.Where("interview = ? AND zju_id = ?", interviewId, zjuId).Delete(&InterviewSchedule{})
}

// 添加指定的面试信息，将检查时间（该面试必须未开始）、冻结和可用容量（使用事务查询、插入）。
//
// 注意不检查指定的学生是否有权限选择面试，也不检查其是否选择过面试，只检查面试能否再加人
func AddScheduledId(interviewInst Interview, zjuId PersonId) (code int, obj any) {
	if interviewInst.StartAt.Before(time.Now()) {
		return utils.Message("面试已开始", 400, 32)
	}
	switch interviewInst.Status {
	case Frozen:
		return utils.Message("面试已冻结", 400, 33)
	case Auto:
		if db.Transaction(func(tx *gorm.DB) error {
			var usedCapacity int64
			tx.
				Model(&InterviewSchedule{}).
				Where("interview = ?", interviewInst.Id).
				Count(&usedCapacity)
			if usedCapacity >= int64(interviewInst.Capacity) {
				return errors.New("面试已满")
			}
			tx.Create(InterviewSchedule{
				ZjuId:     zjuId,
				Interview: interviewInst.Id,
			})
			return nil
		}) != nil {
			return utils.Message("面试已满", 400, 34)
		} else {
			return utils.Success()
		}
	case UnlimitedCapacity:
		db.Create(InterviewSchedule{
			ZjuId:     zjuId,
			Interview: interviewInst.Id,
		})
		return utils.Success()
	default:
		return utils.Message("面试状态未知", 500, 32)
	}
}

// 查看对于指定志愿是否安排了面试
func GetInterviewByIntent(formId uint32, zjuId PersonId, depart uint32, step StepType) *Interview {
	var iv Interview
	result := db.
		Model(&Interview{}).
		Where("interviews.form = ? AND interviews.depart = ? AND interviews.step = ?", formId, depart, step).
		Joins("JOIN interview_schedules ON interviews.id = interview_schedules.interview AND interview_schedules.zju_id = ?", zjuId).
		First(&iv)
	if result.Error != nil || iv.Id <= 0 { //没有找到
		return nil
	}
	return &iv
}
