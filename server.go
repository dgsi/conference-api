package main

import (
	"os"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	h "conference/dgsi/api/handlers"
	m "conference/dgsi/api/models"
	"conference/dgsi/api/config"
	"github.com/jinzhu/gorm"
	"github.com/gin-gonic/contrib/jwt"
)

func main() {
	db := *InitDB()
	router := gin.Default()
	LoadAPIRoutes(router, &db)
}

func LoadAPIRoutes(r *gin.Engine, db *gorm.DB) {
	private := r.Group("/api/v1")
	public := r.Group("/api/v1")
	private.Use(jwt.Auth(config.GetString("TOKEN_KEY")))

	//manage users
	userHandler := h.NewUserHandler(db)
	public.GET("/users", userHandler.Index)
	public.PUT("/users/:user_id", userHandler.Update)
	public.POST("/user", userHandler.Create)
	public.POST("/login", userHandler.Login)

	//manage rooms
	roomHandler := h.NewRoomHandler(db)
	public.GET("/rooms", roomHandler.Index)
	public.GET("/rooms/:id", roomHandler.Show)
	public.PUT("/rooms/:id", roomHandler.Update)
	public.POST("/room", roomHandler.Create)

	//manage topics
	topicHandler := h.NewTopicHandler(db)
	public.POST("/topic", topicHandler.Create)
	public.GET("/topics", topicHandler.Index)
	public.GET("/topics/:topic_id", topicHandler.Show)
	public.GET("/room/:room_id/topics", topicHandler.RoomTopics)
	public.PUT("/topic/:topic_id", topicHandler.Update)
	public.DELETE("/topic/delete/:topic_id", topicHandler.Delete)

	//manage members
	memberHandler := h.NewMemberHandler(db)
	public.POST("/member", memberHandler.Create)
	public.GET("/members", memberHandler.Index)
	public.GET("/members/:member_id", memberHandler.Show)

	//manage attendance
	attendanceHandler := h.NewAttendanceHandler(db)
	public.POST("/attendance", attendanceHandler.Create)
	public.GET("/attendees/room/:room_id", attendanceHandler.AttendeesByRoom)

	//manage room assignments
	assignmentHandler := h.NewRoomAssignmentHandler(db)
	public.GET("/assignments", assignmentHandler.Index)
	public.POST("/room/assign", assignmentHandler.Create)
	public.GET("/assignments/user/:user_id", assignmentHandler.GetAssignementByUser)
	public.GET("/assignments/room/:room_id", assignmentHandler.GetAssigneePerRoom)
	public.DELETE("/assignments/delete/:id", assignmentHandler.Delete)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("PORT ---> ",port)
	r.Run(fmt.Sprintf(":%s", port))
}

func InitDB() *gorm.DB {
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.GetString("DB_USER"), config.GetString("DB_PASS"),
		config.GetString("DB_HOST"), config.GetString("DB_PORT"),
		config.GetString("DB_NAME"))
	log.Printf("\nDatabase URL: %s\n", dbURL)

	_db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to the database:  %s", err))
	}
	_db.DB()
	_db.LogMode(true)
	_db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&m.User{},&m.Room{},&m.Topic{},&m.Member{},&m.Attendance{},&m.RoomAssignment{})
	return _db
}