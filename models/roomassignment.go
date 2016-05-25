package models

type RoomAssignment struct {
	BaseModel
	UserId int `json:"user_id" form:"user_id" binding:"required"`
	RoomId int `json:"room_id" form:"room_id" binding:"required"`
}