package models

import (
	"time"
)

type Topic struct {
	BaseModel
	Title string `json:"title" form:"title" binding:"required"`
	Speaker string `json:"speaker" form:"speaker" binding:"required"`
	RoomNo int `json:"room_no" form:"room_no" binding:"required"`
	StartTime time.Time `json:"start_time" form:"start_time" binding:"required"`
	EndTime time.Time `json:"end_time" form:"end_time" binding:"required"`
}	
