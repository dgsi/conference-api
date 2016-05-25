package models

type QryUser struct {
	UserId int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	UserRole string `json:"user_role"`
	Username string `json:"username"`
	Status string `json:"user_status"`
	Password string `json:"-"`
	IsDefaultPassword bool `json:"is_default_password"`
	RoomId int `json:"room_id"`
	RoomNo string `json:"room_no"`
	Capacity int `json:"capacity"`
}