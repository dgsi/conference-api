package handlers

import(
	"net/http"
	"time"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/gin"
	m "conference/dgsi/api/models"
)

type MemberHandler struct {
	db *gorm.DB
}

func NewMemberHandler(db *gorm.DB) *MemberHandler {
	return &MemberHandler{db}
}

func (handler MemberHandler) Create (c *gin.Context) {
	var newMember m.Member
	c.Bind(&newMember)

	member := m.Member{}
	query := handler.db.Last(&member)
	var customId string
	if query.RowsAffected > 0 {
		lastId,_ := strconv.Atoi(member.CustomId)
		customId = strconv.Itoa(lastId+1)
	} else {
		year := strconv.Itoa(time.Now().Year())
		customId = year + "000001"
	}
	
	newMember.CustomId = customId
	result := handler.db.Create(&newMember)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusCreated, newMember)
	} else {
		respond(http.StatusBadRequest, result.Error.Error(),c,true)
	}
}