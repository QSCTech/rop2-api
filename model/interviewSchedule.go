package model

import "time"

type InterviewSchedule struct {
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"` //面试安排id

	ZjuId     PersonId `json:"zjuId" gorm:"type:char(10);uniqueIndex:uni_zjuId_interview"` //学号
	Interview uint32   `json:"interview" gorm:"not null;uniqueIndex:uni_zjuId_interview"`  //面试id

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
