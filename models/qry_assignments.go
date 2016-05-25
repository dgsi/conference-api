package models

type QryAssignment struct {
	AssignmentId int `json:"assignment_id"`
	RoomId int `json:"room_id"`
	RoomNo string `json:"room_no"`
	UserId int `json:"user_id"`
	User string `json:"user"`
	Username string `json:"username"`
}