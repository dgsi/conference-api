package handlers

import (
	"fmt"
	"net/http"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	m "conference/dgsi/api/models"
)

type TopicHandler struct {
	db *gorm.DB
}

func NewTopicHandler(db *gorm.DB) *TopicHandler {
	return &TopicHandler{db}
}

//create new topic
func (handler TopicHandler) Create(c *gin.Context) {
	var newTopic m.Topic
	c.BindJSON(&newTopic)

	if newTopic.EndTime.Before(newTopic.StartTime) {
		respond(http.StatusBadRequest,"End time must not be earlier than start time",c,true)
	} else if newTopic.StartTime.Hour() == newTopic.EndTime.Hour() {
		respond(http.StatusBadRequest,"Invalid start time and end time",c,true)
	} else {
		topics := []m.Topic{}
		conflictTopic := m.Topic{}
		proceedWithSaving := true
		query := handler.db.Where("room_no = ?",newTopic.RoomNo).First(&topics)
		if query.RowsAffected == 1 {
			for _,t := range topics {
				hour,min,sec := t.StartTime.Clock()
				fmt.Printf("\n hour --> %v min --> %v sec --> %v", hour,min,sec)
				if (newTopic.StartTime.Hour() >= t.StartTime.Hour()) && (newTopic.StartTime.Hour() <= t.EndTime.Hour()) {
				 	proceedWithSaving = false
				 	conflictTopic = t
				 	break
				}
			}
		} 

		if proceedWithSaving {
			result := handler.db.Create(&newTopic)

			if result.RowsAffected == 1 {
				c.JSON(http.StatusCreated, newTopic)
			} else {
				respond(http.StatusBadRequest, result.Error.Error(),c,true)
			}		
		} else {
			respond(http.StatusBadRequest,fmt.Sprintf("Sorry but your desired schedule is already taken by Speaker %s",conflictTopic.Speaker),c,true)
		}
		
	}
}