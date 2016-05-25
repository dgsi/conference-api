package handlers

import (
	"time"
	"fmt"
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
	handler.db.Where("room_no = ?",topic_id).Find(&topics)
	c.JSON(http.StatusOK,topics)
}

//create new topic
func (handler TopicHandler) Create(c *gin.Context) {
	var newTopic m.Topic
	c.Bind(&newTopic)

	fmt.Printf("\nBEFORE PARSE --> %v",c.PostForm("start_time"))
	loc,_ := time.LoadLocation("Asia/Manila")
	startDate,_ := time.ParseInLocation(time.RFC3339,c.PostForm("start_time"),loc)
	endDate,_ := time.ParseInLocation(time.RFC3339,c.PostForm("end_time"),loc)

	fmt.Printf("\nAFTER PARSE ---> %v\n\n",startDate)
	newTopic.StartTime = startDate
	newTopic.EndTime = endDate

	if newTopic.EndTime.Before(newTopic.StartTime) {
		respond(http.StatusBadRequest,"End time must not be earlier than start time",c,true)
	} else if newTopic.StartTime.Hour() == newTopic.EndTime.Hour() {
		respond(http.StatusBadRequest,"Invalid start time and end time",c,true)
	} else {
		fmt.Println("1")
		topics := []m.Topic{}
		query := handler.db.Where("room_no = ? AND (start_time <= ? AND end_time > ?)",newTopic.RoomNo,newTopic.StartTime,newTopic.StartTime).Find(&topics)
		fmt.Println("2")
		if query.RowsAffected > 0 {
			respond(http.StatusBadRequest,"Unable to create topic found conflict with other schedules, Please double check the topic schedules",c,true)
		} else {
			result := handler.db.Create(&newTopic)
			fmt.Printf("\nAFTER SAVING --> %v",newTopic.StartTime)
			if result.RowsAffected > 0 {
				c.JSON(http.StatusCreated, newTopic)
			} else {
				respond(http.StatusBadRequest,result.Error.Error(),c,true)
			}
		}	
	}
}

//update schedule
func (handler TopicHandler) Update(c *gin.Context) {
	topic_id := c.Param("topic_id")
	topic := m.Topic{}

	roomNo,_ := strconv.Atoi(c.PostForm("room_no"))
	room := m.Room{}
	roomQuery := handler.db.Where("id = ?",roomNo).First(&room)

	if roomQuery.RowsAffected > 0 {
		query := handler.db.Where("id = ?",topic_id).First(&topic)
		if query.RowsAffected > 0 {
			canUpdate := true

			loc,_ := time.LoadLocation("Asia/Manila")
			startTime,_ := time.ParseInLocation(time.RFC3339,c.PostForm("start_time"),loc)
			endTime,_ := time.ParseInLocation(time.RFC3339,c.PostForm("end_time"),loc)

			//validation for start time
			if c.PostForm("start_time") != "" {

				if endTime.Before(startTime) {
					canUpdate = false
					respond(http.StatusBadRequest,"End time must not be earlier than start time 111",c,true)
				} else if startTime.Hour() == topic.EndTime.Hour() {
					canUpdate = false
					respond(http.StatusBadRequest,"Invalid start time and end time",c,true)
				} else {
					topics := []m.Topic{}
					conflict := handler.db.Where("id != ? AND room_no = ? AND (start_time <= ? AND end_time > ?)",topic_id,roomNo,startTime,startTime).Find(&topics)

					if conflict.RowsAffected > 0 {
						canUpdate = false
						respond(http.StatusBadRequest,"Unable to update topic found conflict with other schedules, Please double check the topic schedules",c,true)
					} else {
						topic.StartTime = startTime
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

				topic.StartTime = startTime
				topic.EndTime = endTime

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

func (handler TopicHandler) Delete(c *gin.Context) {
	topic_id := c.Param("topic_id")	
	attendance := m.Attendance{}
	query := handler.db.Where("topic_id = ?",topic_id).First(&attendance)
	if query.RowsAffected > 0 {
		respond(http.StatusBadRequest,"This topic already have attendees, cannot proceed with deletion",c,true)
	} else {
		result := handler.db.Where("id = ?",topic_id).Delete(m.Topic{})
		if result.RowsAffected > 0 {
			respond(http.StatusOK,"Topic successfully deleted",c,false)
		} else {
			respond(http.StatusBadRequest,result.Error.Error(),c,true)
		}
	}
}

