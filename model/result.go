package model

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 报名者答卷。同一人报名多个志愿部门，仅一份答卷，但生成多个Intent。
type Result struct {
	Form  uint32 `json:"form" gorm:"primaryKey"`
	ZjuId string `json:"zjuId" gorm:"type:char(10);primaryKey"`
	//实际答卷内容
	Content string `json:"content" gorm:"not null;type:json"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}

// 注意：重新提交答卷/被管理员重置“已填表”状态，都不会影响已报名的面试
// 即：即使不在对应阶段，候选人可能仍然在某场面试的名单中(类似面试后被拒绝，但仍然在原来的面试名单中)
func SaveFullResult(formId uint32, zjuId PersonId, phone string, content string, intentDeparts []uint32) error {
	//开事务，保证一致性
	return db.Transaction(func(tx *gorm.DB) error {
		//更新个人信息。person表中一定存在对应的学号(登录时创建)
		if err := tx.Model(&Person{}).Where("zju_id = ?", zjuId).Update("Phone", phone).Error; err != nil {
			return err
		}

		//更新答卷
		if err := tx.
			//id由数据库自动生成；create_at和update_at由gorm添加
			Select("form", "zju_id", "content").
			Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns([]string{"content", "update_at"})}).
			Create(&Result{
				Form:    formId,
				ZjuId:   zjuId,
				Content: content,
			}).Error; err != nil {
			return err
		}

		//更新志愿部门
		if err := tx. //删除不再存在的志愿
				Where("form = ? AND zju_id = ?", formId, zjuId).
				Where("depart NOT IN ?", intentDeparts).
				Delete(&Intent{}).Error; err != nil {
			return err
		}
		var interviewScheduleIdsToDelete []uint32
		tx. //查询需要删除的面试安排
			Model(&InterviewSchedule{}).
			Joins("INNER JOIN interviews on interview_schedules.interview = interviews.id").
			Where("form = ? AND zju_id = ?", formId, zjuId).
			Where("interviews.depart NOT IN ?", intentDeparts).
			Pluck("interview_schedules.id", &interviewScheduleIdsToDelete)
		if len(interviewScheduleIdsToDelete) > 0 {
			if err := tx. //实际删除面试安排
					Delete(&InterviewSchedule{}, "id IN ?", interviewScheduleIdsToDelete).Error; err != nil {
				return err
			}
		}

		intents := make([]*Intent, len(intentDeparts))
		for i, v := range intentDeparts {
			intents[i] = &Intent{
				Form:   formId,
				ZjuId:  zjuId,
				Depart: v,
				Order:  int8(i + 1),
			}
		}
		if err := tx.
			//id由数据库自动生成；step默认为0；=
			//create_at和update_at由gorm添加
			Select("Form", "ZjuId", "Depart", "Order").
			Clauses(clause.OnConflict{ //冲突时仅更新order和update_at
				DoUpdates: clause.AssignmentColumns([]string{"order", "update_at"})}).
			Create(intents).Error; err != nil {
			return err
		}
		return nil
	})
}

type ResultDetail struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Content string `json:"content"`
}

func GetResult(formId uint32, zjuId string) *ResultDetail {
	var contentArr []string
	db.Model(&Result{}).Where("form = ? AND zju_id = ?", formId, zjuId).Pluck("content", &contentArr)
	if len(contentArr) == 0 {
		return nil
	}
	content := contentArr[0]
	person := FindPerson(zjuId)
	return &ResultDetail{
		Name:    person.Name,
		Phone:   *(person.Phone), //填表必须提交手机号，此处一定不为空指针
		Content: content,
	}
}
