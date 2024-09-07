package model

import (
	"time"

	"gorm.io/gorm"
)

// 候选人的某个志愿
type Intent struct {
	//自增主键，志愿id
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement"`
	//来源表单（也确定了申请的组织）
	Form uint32 `json:"form" gorm:"not null;uniqueIndex:zjuid_form_intent,sort:desc"`
	//申请人学号
	ZjuId string `json:"zjuId" gorm:"type:char(10);not null;uniqueIndex:zjuid_form_intent,sort:desc"`
	//志愿部门，可能为默认部门（如果未选择志愿部门）
	Depart uint32 `json:"depart" gorm:"not null;uniqueIndex:zjuid_form_intent"`
	//志愿排序。1~127=第n志愿
	Order int8 `json:"order" gorm:"not null;uniqueIndex:zjuid_form_intent,sort:asc"`
	//当前所在阶段。1~127=第n阶段(可重命名)
	Step StepType `json:"step" gorm:"not null;default:0"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

func SaveIntents(formId uint32, zjuId string, intentDeparts []uint32) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Intent{}, "form = ? AND zju_id = ?", formId, zjuId).Error; err != nil {
			return err
		}
		intents := make([]Intent, len(intentDeparts))
		for i, v := range intentDeparts {
			intents[i] = Intent{
				Form:   formId,
				ZjuId:  zjuId,
				Depart: v,
				Order:  int8(i + 1),
			}
		}
		return tx.Select("Form", "ZjuId", "Depart", "Order").Create(intents).Error
	})
}

type IntentOutline struct {
	Id     uint32 `json:"id"` //志愿id
	Name   string `json:"name"`
	ZjuId  string `json:"zjuId"`
	Phone  string `json:"phone"`
	Depart uint32 `json:"depart"`
	Order  int8   `json:"order"`

	//候选人报名的面试信息（可能为空）
	InterviewTime *time.Time `json:"interviewTime"`
}
type IntentList struct {
	Intents       []IntentOutline `json:"intents"`
	Count         int64           `json:"count"`
	FilteredCount int64           `json:"filteredCount"`
}

func ListIntents(formId uint32, departs []uint32, step StepType, offset, limit int, filter string) IntentList {
	var count int64 //在指定部门、阶段的总志愿数
	db.Table("intents").
		Joins("INNER JOIN people ON intents.zju_id = people.zju_id").
		Where("intents.form = ?", formId).
		Where("intents.depart IN ?", departs).
		Where("intents.step = ?", step).
		Count(&count)

	if filter == "" {
		filter = "^" //匹配所有
	}

	intents := make([]IntentOutline, 0)
	db.
		Table("intents").
		Select("people.name, people.zju_id, people.phone, intents.order, intents.depart, intents.id, scheduled_interviews.start_at as InterviewTime").
		//查询每个志愿 候选人的姓名、学号、手机号信息
		//此处INNER JOIN确保intents和people表都有zju_id对应的信息
		Joins("INNER JOIN people ON intents.zju_id = people.zju_id").
		//根据志愿信息查询报名过的面试
		//先把interview_schedules和interviews两表LEFT JOIN生成派生表scheduled_interviews(确保可在派生表中直接查询面试的form、depart、step等信息)
		//然后把intents和scheduled_interviews两表LEFT JOIN，并限制派生表中的zju_id、form、depart、step与intents中对应
		Joins("LEFT JOIN (SELECT interview_schedules.interview, interview_schedules.zju_id, interviews.* from interview_schedules LEFT JOIN interviews ON interview_schedules.interview = interviews.id) `scheduled_interviews` ON scheduled_interviews.zju_id = intents.zju_id AND scheduled_interviews.form = intents.form AND scheduled_interviews.depart = intents.depart AND scheduled_interviews.step = intents.step").
		Where("intents.form = ?", formId).
		Where("intents.depart IN ?", departs).
		Where("intents.step = ?", step).
		Where("people.name REGEXP ? OR people.zju_id REGEXP ? OR people.phone REGEXP ?", filter, filter, filter).
		Order("intents.create_at ASC, intents.order ASC").
		Offset(offset).Limit(limit).
		Scan(&intents)

	var filteredCount int64 //在指定部门、阶段+filter过滤后的志愿数，实际上就是上方的查询改为count，不带offset和limit
	//据称mysql会优化此类连续查询，不建议使用特殊常量获取count
	db.
		Table("intents").
		Joins("INNER JOIN people ON intents.zju_id = people.zju_id").
		Where("intents.form = ?", formId).
		Where("intents.depart IN ?", departs).
		Where("intents.step = ?", step).
		Where("people.name REGEXP ? OR people.zju_id REGEXP ? OR people.phone REGEXP ?", filter, filter, filter).
		Count(&filteredCount)

	return IntentList{Intents: intents, Count: count, FilteredCount: filteredCount}
}

func SetIntents(formId uint32, intentIds []uint32, step StepType) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&Intent{}).Where("form = ? AND id IN ?", formId, intentIds).Update("step", step).Error
	})
}

func QueryIntentsOfPerson(formId uint32, zjuId string) []Intent {
	intents := make([]Intent, 0)
	db.Where("form = ? AND zju_id = ?", formId, zjuId).Order("`order` ASC").Find(&intents)
	return intents
}
