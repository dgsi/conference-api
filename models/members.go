package models

import (
	"errors"
	"strings"
)

type Member struct {
	BaseModel
	CustomId string `json:"custom_id" form:"custom_id" binding:"required"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName string `json:"last_name" form:"last_name" binding:"required"`
	Address string `json:"address" form:"address" binding:"required"`
	ContactNo string `json:"contact_no" form:"contact_no" binding:"required"`
	Email string `json:"email" form:"email" binding:"required"`
	Gender string `json:"gender" form:"gender" binding:"required"`
	Status string `json:"status" form:"status"`
}

func (m *Member) BeforeCreate() (err error) {
	m.Status = "active"
	if strings.TrimSpace(m.FirstName) == "" {
		err = errors.New("Please specify the member's first name")
	} else if strings.TrimSpace(m.LastName) == "" {
		err = errors.New("Please specify the member's last name")
	} else if strings.TrimSpace(m.Gender) == "" {
		err = errors.New("Please specify the member's gender")
	}
	return
}