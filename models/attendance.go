package models

import "time"

type Attendance struct {
	Id int `json:"id" gorm:"AUTO_INCREMENT; primary_key"`
	MemberId string `json:"member_id" form:"member_id" binding:"required"`
	TopicId int `json:"topic_id" form:"topic_id"`	
	RoomId int `json:"room_id" form:"room_id" binding:"required"`
	Status string `json:"status" form:"status"`
	ScannedBy int `json:"scanned_by" form:"scanned_by" binding:"required"`
	Mode string `json:"mode" form:"mode" binding:"required"` 
	TITO time.Time `json:"tito" form:"tito"`
}