//注：相关功能暂时延期

package model

import "time"

//通知模板。
type Template struct {
	//主键，自动递增
	Id uint32 `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`

	//归属的组织（不是部门）
	Owner uint32 `json:"owner" gorm:"not null;uniqueIndex:uni_owner_name"`
	//自定义的标识，便于查找选择，与发送无关，须在组织内唯一
	Name string `json:"name" gorm:"not null;type:varchar(40);uniqueIndex:uni_owner_name"`

	//实际通知内容（可含变量）
	Content string `json:"content" gorm:"not null"`

	CreateAt time.Time `json:"createAt" gorm:"not null;autoCreateTime"`
	UpdateAt time.Time `json:"updateAt" gorm:"not null;autoUpdateTime"`
}
