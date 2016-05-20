package handlers

import(
	"net/http"
	"time"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	m "conference/dgsi/api/models"
)

type AttendanceHandler struct {
	db *gorm.DB
}

func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
	return &AttendanceHandler{db}
}

func (handler AttendanceHandler) Create (c *gin.Context) {
	var newAttendance m.Attendance
	c.Bind(&newAttendance)

	//find topic based on time stamp
	topics := []m.Topic{}
	query := handler.db.Where("room_no = ?",newAttendance.RoomId).Find(&topics)

	if query.RowsAffected > 0 {

		if newAttendance.Mode == "manual" {
			newAttendance.Status = c.PostForm("status")
			ProceedWithSaving(true, newAttendance,topics,handler,c)
		} else {
			ProceedWithSaving(false, newAttendance,topics,handler,c)
		}
	} else {
		respond(http.StatusBadRequest,"There are no scheduled topic in this room right now",c,true)
	}		
}

func ProceedWithSaving(isManual bool, newAttendance m.Attendance, topics []m.Topic, handler AttendanceHandler, c *gin.Context) {
	tito,_ := time.Parse(time.RFC3339,c.PostForm("tito"))
	var found bool = false

	for _,t := range topics {
		fmt.Printf("\nSTART TIME --> %v",t.StartTime)
		fmt.Printf("\nEND TIME --> %v",t.EndTime)
		fmt.Printf("\nTITO TIME --> %v",tito)
		fmt.Println("\n\n")
		if ((t.StartTime.Equal(tito) || t.StartTime.Before(tito)) && (t.EndTime.Equal(tito) || t.EndTime.After(tito))) {
			newAttendance.TopicId = t.Id
			newAttendance.TITO = tito
			found = true
			break
		} 
	}

	if found {
		var canProceedToSaving bool = true
		var errMsg string
		status := m.Attendance{}
		queryStatus := handler.db.Where("room_id = ? AND member_id = ?",newAttendance.RoomId,newAttendance.MemberId).Last(&status)
		
		if queryStatus.RowsAffected > 0 {
			fmt.Printf("\n rows affected greater than 0")
			fmt.Printf("\n mode --> %v",isManual)
			fmt.Printf("\n status.RoomId --> %v", status.RoomId)
			fmt.Printf("\n newAttendance.RoomdId -> %v", newAttendance.RoomId)
			fmt.Printf("\n c.PostForm(status) --> %v", c.PostForm("status"))
			
			if isManual {
				if status.RoomId == newAttendance.RoomId {
					if status.Status == c.PostForm("status") {
						fmt.Printf("\nEQUAL")
						
						if (status.Status == "time in") {
							errMsg = "Member already clocked in"
						} else {
							errMsg = "Member already clocked out"
						}
						fmt.Printf("\n errMsg --> %v",errMsg)
						fmt.Printf("\n\n")
						canProceedToSaving = false
					} 
				} else {
					errMsg = "asjdlajsdlkajskldja"
					canProceedToSaving = false
				}
			} else {
				if status.RoomId == newAttendance.RoomId {
					if status.Status == "time out" {
						newAttendance.Status = "time in"
					} else {
						newAttendance.Status = "time out"	
					}
				} else {
					canProceedToSaving = false
				}
			}
		} else {
			newAttendance.Status = "time in"
		}

		conflict := m.Attendance{}
		conflictQuery := handler.db.Where("room_id != ? AND member_id = ?",newAttendance.RoomId,newAttendance.MemberId).Last(&conflict)
		
		fmt.Printf("\nCONFLICT QUERY SIZE --> %v",conflictQuery.RowsAffected)

		if (conflictQuery.RowsAffected > 0 && conflict.Status == "time in") {
			errMsg = fmt.Sprintf("You are still not clocking out from room %v",conflict.RoomId)
			canProceedToSaving = false
		} 

		if canProceedToSaving {
			result := handler.db.Create(&newAttendance)

			if result.RowsAffected > 0 {
				c.JSON(http.StatusCreated,newAttendance)
			} else {
				respond(http.StatusBadRequest,result.Error.Error(),c,true)
			}
		} else {
			respond(http.StatusBadRequest,errMsg,c,true)
		}
	} else {
		respond(http.StatusBadRequest,"There are no scheduled topic in this room right now 22",c,true)
	}
}

func (handler AttendanceHandler) AttendeesByRoom(c *gin.Context) {
	room_id := c.Param("room_id")
	attendees := []m.QryAttendance{}
	handler.db.Where("room_id = ? AND status = ?",room_id,"time in").Find(&attendees)
	c.JSON(http.StatusOK, attendees)
}





