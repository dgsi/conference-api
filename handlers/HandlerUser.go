package handlers

import(
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"conference/dgsi/api/config"
	m "conference/dgsi/api/models"
	jwt_lib "github.com/dgrijalva/jwt-go"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db}
}

func (handler UserHandler) Create(c *gin.Context) {
	var newUser m.User
	c.Bind(&newUser)

	//generate auto username
	user := m.User{}	
	handler.db.Last(&user)

	if user.UserName == "" {
		newUser.Increment = "1"
	} else {
		i,_ := strconv.Atoi(newUser.Increment)
		newUser.Increment = strconv.Itoa(i+1)	
	}

	newUser.UserName = newUser.UserRole + newUser.Increment
	result := handler.db.Create(&newUser)

	if result.RowsAffected == 1 {
		//generate jwt
		newUser.Token = generateJWT(newUser.UserName)
		c.JSON(http.StatusCreated, newUser)
	} else {
		respond(http.StatusBadRequest, result.Error.Error(),c,true)
	}
}

func generateJWT(username string) string {
	// Create the token
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
	// Set some claims
	token.Claims["ID"] = username
	token.Claims["exp"] = 0
	// Sign and get the complete encoded token as a string
	tokenString, _ := token.SignedString([]byte(config.GetString("TOKEN_KEY")))
    return tokenString
}