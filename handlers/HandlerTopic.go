package handlers

import (
	"time"
	"strconv"
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
		query := handler.db.Where("room_no = ? AND start_time BETWEEN ? AND ?",newTopic.RoomNo,newTopic.StartTime,newTopic.EndTime).Find(&topics)

		if query.RowsAffected > 0 {
			respond(http.StatusBadRequest,"Unable to create topic found conflict with other schedules, Please double check the topic schedules",c,true)
		} else {
			result := handler.db.Create(&newTopic)

			if result.RowsAffected > 0 {
				c.JSON(http.StatusCreated, newTopic)
			} else {
				respond(http.StatusBadRequest, result.Error.Error(),c,true)
			}
		}	
	}
}

//update schedule
func (handler TopicHandler) Update (c *gin.Context) {
	topic_id := c.Param("topic_id")
	topic := m.Topic{}

	roomNo,_ := strconv.Atoi(c.PostForm("room_no"))
	room := m.Room{}
	roomQuery := handler.db.Where("id = ?",roomNo).First(&room)

	if roomQuery.RowsAffected > 0 {
		query := handler.db.Where("id = ?",topic_id).First(&topic)
		if query.RowsAffected > 0 {
			canUpdate := true

			startTime,_ := time.Parse(time.RFC3339,c.PostForm("start_time"))
			endTime,_ := time.Parse(time.RFC3339,c.PostForm("end_time"))

			//validation for start time
			if c.PostForm("start_time") != "" {

				if endTime.Before(startTime) {
					canUpdate = false
					respond(http.StatusBadRequest,"End time must not be earlier than start time 111",c,true)
				} else if startTime.Hour() == topic.EndTime.Hour() {
					respond(http.StatusBadRequest,"Invalid start time and end time",c,true)
				} else {
					topics := []m.Topic{}
					conflict := handler.db.Where("id != ? AND room_no = ? AND start_time BETWEEN ? AND ?",topic.Id,roomNo,startTime,topic.EndTime).Find(&topics)

					if conflict.RowsAffected > 0 {
						respond(http.StatusBadRequest,"Unable to update topic found conflict with other schedules, Please double check the topic schedules",c,true)
					} else {
						topic.StartTime = startTime
					}	
				}
			}	

			//validation for end time
			if c.PostForm("end_time") != "" {

				if endTime.Before(startTime) {
					canUpdate = false
					respond(http.StatusBadRequest,"End time must not be earlier than start time",c,true)
				} else if endTime.Hour() == topic.StartTime.Hour() {
					respond(http.StatusBadRequest,"Invalid start time and end time",c,true)
				} else {
					topics := []m.Topic{}
					conflict := handler.db.Where("id != ? AND room_no = ? AND start_time BETWEEN ? AND ?",topic.Id,roomNo,topic.StartTime,endTime).Find(&topics)

					if conflict.RowsAffected > 0 {
						respond(http.StatusBadRequest,"Unable to update topic found conflict with other schedules, Please double check the topic schedules",c,true)
					} else {
						topic.EndTime = endTime
					}	
				}
			}	

			if canUpdate {
				//check for updates in topic title
				if c.PostForm("title") != "" {
					topic.Title = c.PostForm("title")
				}
				//check for udpates in speaker
				if c.PostForm("speaker") != "" {
					topic.Speaker = c.PostForm("speaker")
				}
				//check for updates in description
				if c.PostForm("description") != "" {
					topic.Description = c.PostForm("description")
				}
				//check for updates in room no
				if c.PostForm("room_no") != "" {
					topic.RoomNo = roomNo
				}

				update := handler.db.Save(&topic)
				if update.RowsAffected > 0 {
					c.JSON(http.StatusOK,topic)
				} else {
					respond(http.StatusBadRequest,update.Error.Error(),c,true)
				}
			}
		} else {
			respond(http.StatusBadRequest,"Unable to find topic",c,true)
		}
	} else {
		respond(http.StatusBadRequest,"Unable to find room",c,true)
	}
}

