package models

import "time"

type QryAttendance struct {
	Id int `json:"id"`
	Status string `json:"status"`
	TITO time.Time `json:"time_in_time_out"`
	CustomId string `json:"member_id"`
	FirstName string `json:"member_first_name"`
	LastName string `json:"member_last_name"`
	Address string `json:"member_address"`
	ContactNo string `json:"member_contact_no"`
	Email string `json:"email"`
	Gender string `json:"gender"`
	MemberStatus string `json:"member_status"`
	Title string `json:"title"`
	Description string `json:"description"`
	Speaker string `json:"speaker"`
	StartTime time.Time `json:"topic_start_time"`
	EndTime time.Time `json:"topic_end_time"`
	RoomNo string `json:"room_no"`
	RoomStatus string `json:"room_status"`
	Capacity int `json:"room_capacity"`
	ScannedBy string `json:"scanned_by"`
}
