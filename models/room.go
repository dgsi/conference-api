package models

import (
	"errors"
	"strings"
)

type Room struct {
	BaseModel
	RoomNo string `json:"room_no" form:"room_no" binding:"required"`
	Status string `json:"status" form:"status"`
	Capacity string `json:"capacity" form:"capacity" binding:"required"`
}

func (r *Room) BeforeCreate() (err error) {
	r.Status = "active"
	if strings.TrimSpace(r.RoomNo) == "" {
		err = errors.New("Please specify the room no")
	} else if strings.TrimSpace(r.Capacity) == "" {
		err = errors.New("Please specify the room's capacity")
	}
	return nil
}