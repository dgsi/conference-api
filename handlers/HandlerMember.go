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

func (handler MemberHandler) Index (c *gin.Context) {
	members := []m.Member{}
	handler.db.Where("status = ? ","active").Order("created_at desc").Find(&members)
	c.JSON(http.StatusOK,members)
}

func (handler MemberHandler) Show (c *gin.Context) {
	member_id := c.Param("member_id")
	member := m.Member{}
	handler.db.Where("status = ? AND custom_id = ?","active",member_id).First(&member)
	c.JSON(http.StatusOK,member)
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