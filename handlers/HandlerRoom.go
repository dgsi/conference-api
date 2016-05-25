package handlers

import(
	"fmt"
	"net/http"
	"strings"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "conference/dgsi/api/models"
)

type RoomHandler struct {
	db *gorm.DB
}

func NewRoomHandler(db *gorm.DB) *RoomHandler {
	return &RoomHandler{db}
}

//show all rooms
func (handler RoomHandler) Index(c *gin.Context) {
	rooms := []m.Room{}
	handler.db.Where("status = ?","active").Order("created_at desc").Find(&rooms)
	c.JSON(http.StatusOK,rooms)
}

//create new room
func (handler RoomHandler) Create(c *gin.Context) {
	var newRoom m.Room
	c.Bind(&newRoom)

	//generate auto username
	room := m.Room{}	
	query := handler.db.Last(&room)

	if query.RowsAffected == 0 {
		newRoom.RoomNo = "Room 1"
	} else {
		newRoom.RoomNo = "Room " + strconv.Itoa(room.Id+1)
	}

	//insert record to database
	result := handler.db.Create(&newRoom)
	fmt.Printf("\nrows affected --> %v",result.RowsAffected)

	if result.RowsAffected == 1 {
		c.JSON(http.StatusCreated, newRoom)
	} else {
		respond(http.StatusBadRequest, result.Error.Error(),c,true)
	}
}

//show specific room by id
func (handler RoomHandler) Show(c *gin.Context) {
	id := c.Param("id")
	room := m.Room{}
	query := handler.db.Where("id = ?",id).First(&room)
	if query.RowsAffected == 1 {
		c.JSON(http.StatusOK,room)
	} else {
		respond(http.StatusBadRequest,"Room record not found!",c,true)
	}
}

func (handler RoomHandler) Update(c *gin.Context) {
	id := c.Param("id")
	capacity := c.PostForm("capacity")
	room := m.Room{}
	query := handler.db.Where("id = ?",id).First(&room)
	if query.RowsAffected == 1 {
		if strings.TrimSpace(capacity) == "" {
			respond(http.StatusBadRequest,"Please specify the new room's capacity",c,true)
		} else {
			newCapacity,_ := strconv.Atoi(capacity)
			if newCapacity < 10 {
				respond(http.StatusBadRequest,"Room capacity must be atleast greater than 10",c,true)
			} else {
				room.Capacity = newCapacity
				status := c.PostForm("status")
				if (strings.TrimSpace("status") != "") {
					room.Status = status
				}
				handler.db.Save(&room)
				c.JSON(http.StatusOK,room)	
			}
		}
	} else {
		respond(http.StatusBadRequest,"Room record not found!",c,true)
	}
}