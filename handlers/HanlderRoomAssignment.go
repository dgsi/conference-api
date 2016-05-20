package handlers

import(
	"net/http"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	m "conference/dgsi/api/models"
)

type RoomAssignmentHandler struct {
	db *gorm.DB
}

func NewRoomAssignmentHandler(db *gorm.DB) *RoomAssignmentHandler {
	return &RoomAssignmentHandler{db}
}

func (handler RoomAssignmentHandler) Index (c *gin.Context) {
	assignments := []m.RoomAssignment{}
	handler.db.Find(&assignments)
	c.JSON(http.StatusOK,assignments)
}

func (handler RoomAssignmentHandler) GetAssignementByUser (c *gin.Context) {
	user_id := c.Param("user_id")
	assignments := []m.RoomAssignment{}
	handler.db.Where("user_id = ?",user_id).Find(&assignments)
	c.JSON(http.StatusOK,assignments)
}

func (handler RoomAssignmentHandler) GetAssigneePerRoom (c *gin.Context) {
	room_id := c.Param("room_id")
	assignments := []m.RoomAssignment{}
	handler.db.Where("room_id = ?",room_id).Find(&assignments)
	c.JSON(http.StatusOK,assignments)
}

func (handler RoomAssignmentHandler) Create (c *gin.Context) {
	var newAssignment m.RoomAssignment
	c.Bind(&newAssignment)

	//check if user is existing
	user := m.User{}
	userQuery := handler.db.Where("id = ?",newAssignment.UserId).First(&user)

	if userQuery.RowsAffected > 0 {
		room := m.Room{}
		roomQuery := handler.db.Where("id = ?",newAssignment.RoomId).First(&room)
		if roomQuery.RowsAffected > 0 {
			//check if user has no room assignment yet
			assignment := m.RoomAssignment{}
			query := handler.db.Where("user_id = ?",newAssignment.UserId).Last(&assignment)

			if query.RowsAffected > 0 {
				respond(http.StatusBadRequest,fmt.Sprintf("Sorry but user is already assigned at room no %v",assignment.RoomId),c,true)
			} else {
				result := handler.db.Create(&newAssignment)
				if result.RowsAffected > 0 {
					c.JSON(http.StatusCreated,newAssignment)
				} else {
					respond(http.StatusBadRequest,result.Error.Error(),c,true)
				}
			}
		} else {
			respond(http.StatusBadRequest,"Room not found",c,true)
		}
	} else {
		respond(http.StatusBadRequest,"User record not found",c,true)
	}
}