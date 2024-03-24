package model

import (
	"time"

	"gorm.io/gorm"
)

type Org struct {
	Id   uint32 `json:"id" gorm:"primaryKey;autoIncrement;<-:false"` //主键，自动递增
	Name string `json:"name" gorm:"type:varchar(80);not null;unique"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`

	DefaultDepart uint32 `json:"defaultDepart" gorm:"uniqueIndex"`
}

func GetOrg(id uint32) *Org {
	var pobj = &Org{}
	result := db.First(pobj, id)
	if result.Error != nil {
		return nil
	} else {
		return pobj
	}
}

// 内部方法。初始化一个组织（包括默认部门、管理员、默认拒信）。
func InitNewOrg(name string, adminZjuId string, adminNickname string) {
	transactionFunc := func(tx *gorm.DB) (err error) {
		org := &Org{
			Name: name,
		}
		if err = tx.Select("Name").Create(org).Error; err != nil {
			return
		}
		defaultDepart := &Depart{
			Name:  "默认部门",
			Owner: org.Id,
		}
		if err = tx.Select("Name", "Owner").Create(defaultDepart).Error; err != nil {
			return
		}
		org.DefaultDepart = defaultDepart.Id
		if err = tx.Select("DefaultDepart").Save(org).Error; err != nil {
			return
		}
		admin := &Admin{
			ZjuId:    adminZjuId,
			At:       org.Id,
			Nickname: adminNickname,
			Level:    Maintainer,
		}
		if err = tx.Select("ZjuId", "At", "Nickname", "Level").Create(admin).Error; err != nil {
			return
		}
		//一些默认文本可能需要修改
		defaultRejectTemplate := &Template{
			Owner:   org.Id,
			Name:    "默认拒信",
			Content: "很遗憾，您未能成功加入{组织}。",
		}
		if err = tx.Select("Owner", "Name", "Content").Create(defaultRejectTemplate).Error; err != nil {
			return
		}
		defaultRejectStage := &Stage{
			Owner:   org.DefaultDepart,
			Step:    Rejected,
			OnEnter: &defaultRejectTemplate.Id,
		}
		if err = tx.Select("Owner", "Step", "OnEnter").Create(defaultRejectStage).Error; err != nil {
			return
		}

		//成功完成，没有错误
		return
	}
	db.Transaction(func(tx *gorm.DB) error {
		err := transactionFunc(tx)
		//或许可以处理一下err
		return err
	})
}
