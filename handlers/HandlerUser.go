package handlers

import(
	"net/http"
	"strconv"
	"strings"
	"fmt"
	"crypto/aes"
    "crypto/cipher"
    "encoding/base64"
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

//get listing of all active user
func (handler UserHandler) Index(c *gin.Context) {
	users := []m.User{}
	handler.db.Where("status = ?","active").Find(&users)
	c.JSON(http.StatusOK,users)
}

//create new user
func (handler UserHandler) Create(c *gin.Context) {
	var newUser m.User
	c.Bind(&newUser)

	//generate auto username
	user := m.User{}	
	query := handler.db.Last(&user)

	if query.RowsAffected == 0 {
		newUser.Increment = "1"
	} else {
		i,_ := strconv.Atoi(newUser.Increment)
		newUser.Increment = strconv.Itoa(i+1)	
	}

	newUser.Username = newUser.UserRole + newUser.Increment
	result := handler.db.Create(&newUser)

	if result.RowsAffected == 1 {
		//generate jwt
		newUser.Token = generateJWT(newUser.Username)
		c.JSON(http.StatusCreated, newUser)
	} else {
		respond(http.StatusBadRequest, result.Error.Error(),c,true)
	}
}

func (handler UserHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.TrimSpace(username) == "" {
		respond(http.StatusBadRequest,"Username is required",c,true)
	} else if strings.TrimSpace(password) == "" {
		respond(http.StatusBadRequest,"Password is required",c,true)
	} else {
		
		user := m.User{}	
		result := handler.db.Where("username = ?", username).Find(&user)
		if result.RowsAffected == 1 {
			decryptedPassword := decrypt([]byte(config.GetString("CRYPT_KEY")), user.Password)
			if password == decryptedPassword {
				user.Token = generateJWT(user.Username)
				c.JSON(http.StatusOK, user)
			} else {
				respond(http.StatusBadRequest,"Account not found!",c,true)
			}
		} else {
			respond(http.StatusBadRequest,"Account not found!",c,true)
		}
	}	
}

//generate java web token
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

// decrypt from base64 to decrypted string
func decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}