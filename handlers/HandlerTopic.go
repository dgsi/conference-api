package handlers

import (
	"time"
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

//get all topics
func (handler TopicHandler) Index(c *gin.Context) {
	topics := []m.Topic{}
	handler.db.Find(&topics)
	c.JSON(http.StatusOK,topics)
}

//show topic by topic id
func (handler TopicHandler) Show(c *gin.Context) {
	topic_id := c.Param("topic_id")
	topic := m.Topic{}
	query := handler.db.Where("id = ?",topic_id).First(&topic)
	if query.RowsAffected > 0 {
		c.JSON(http.StatusOK,topic)
	} else {
		respond(http.StatusBadRequest,"Unable to find room",c,true)
	}
}

//show topics in a room
func (handler TopicHandler) RoomTopics(c *gin.Context) {
	topic_id := c.Param("room_id")
	topics := []m.Topic{}
	query := handler.db.Where("room_no = ?",topic_id).First(&topics)
	if query.RowsAffected > 0 {
		c.JSON(http.StatusOK,topics)
	} else {
		respond(http.StatusBadRequest,"Unable to find topic",c,true)
	}
}


//create new topic
func (handler TopicHandler) Create(c *gin.Context) {
	var newTopic m.Topic
	c.Bind(&newTopic)

	startDate,_ := time.Parse(time.RFC3339,c.PostForm("start_time"))
	endDate,_ := time.Parse(time.RFC3339,c.PostForm("end_time"))

	newTopic.StartTime = startDate
	newTopic.EndTime = endDate

	if newTopic.EndTime.Before(newTopic.StartTime) {
		respond(http.StatusBadRequest,"End time must not be earlier than start time",c,true)
	} else if newTopic.StartTime.Hour() == newTopic.EndTime.Hour() {
		respond(http.StatusBadRequest,"Invalid start time and end time",c,true)
	} else {
		topics := []m.Topic{}
		query := handler.db.Where("start_time BETWEEN ? AND ?",newTopic.StartTime,newTopic.EndTime).Find(&topics)

		if query.RowsAffected == 1 {
			respond(http.StatusBadRequest,"With schedule conflicts",c,true)
		} else {
			result := handler.db.Create(&newTopic)

			if result.RowsAffected == 1 {
				c.JSON(http.StatusCreated, newTopic)
			} else {
				respond(http.StatusBadRequest, result.Error.Error(),c,true)
			}
		}	
	}
}

//update schedule
func (handler TopicHandler) Update (c *gin.Context) {

}

