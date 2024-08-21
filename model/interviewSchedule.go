package model

import "time"

type InterviewSchedule struct {
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"` //面试安排id

	ZjuId     PersonId `json:"zjuId" gorm:"type:char(10);uniqueIndex:uni_zjuId_interview"` //学号
	Interview uint32   `json:"interview" gorm:"not null;uniqueIndex:uni_zjuId_interview"`  //面试id，同样指定了此次面试安排归属组织、表单、部门、阶段

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

func GetScheduledIds(interviewId uint32) []PersonId {
	result := make([]PersonId, 0)
	db.Model(&InterviewSchedule{}).Where("interview = ?", interviewId).Pluck("zju_id", &result)
	return result
}

func RemoveScheduledId(interviewId uint32, zjuId PersonId) {
	db.Where("interview = ? AND zju_id = ?", interviewId, zjuId).Delete(&InterviewSchedule{})
}

func AddScheduledId(interviewId uint32, zjuId PersonId) {
	db.Create(&InterviewSchedule{
		ZjuId:     zjuId,
		Interview: interviewId,
	})
}

//查看对于指定志愿是否安排了面试
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
