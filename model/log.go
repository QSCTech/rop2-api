package model

import "time"

type Log struct {
	Id    uint32    `json:"id" gorm:"primaryKey;autoIncrement;<-:false"`
	Time  time.Time `json:"time" gorm:"type:timestamp;not null;autoCreateTime;default:now()"`
	ZjuId PersonId  `json:"zjuId" gorm:"type:char(10);not null"`

	Method string  `json:"method" gorm:"type:char(10);not null"` //请求方法
	Path   string  `json:"path" gorm:"not null"`                 //请求路径(包括query)，不包括域名
	Body   *[]byte `json:"body" gorm:"type:MediumBlob"`          //请求body
	Status int16   `json:"status" gorm:"not null"`               //响应的HTTP状态码(非自定义状态码)
}

func CreateLog(zjuId string, method string, path string, body *[]byte, status int16) {
	db.
		Select("zju_id", "method", "path", "body", "status").
		Create(&Log{ZjuId: zjuId, Method: method, Path: path, Body: body, Status: status})
}
