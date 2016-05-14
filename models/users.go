package models

import (
	"errors"
	"strings"
	"io"	
	"crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "conference/dgsi/api/config"
)

type User struct {
	BaseModel
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName string `json:"last_name" form:"last_name" binding:"required"`
	Status string `json:"status" form:"status"`
	UserRole string `json:"user_role" form:"user_role" binding:"required"`
	Username string `json:"username" form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	IsDefaultPassword bool `json:"is_default_password"`
	Increment string
	Token string `json:"token"`
}

func (u *User) BeforeCreate() (err error) {
	if strings.TrimSpace(u.FirstName) == "" {
		err = errors.New("Please specify the user's first name")
	} else if strings.TrimSpace(u.LastName) == "" {
		err = errors.New("Please specify the user's last name")
	} else if strings.TrimSpace(u.UserRole) == "" {
		err = errors.New("Please specify the user's role")
	} 
	//set default status 
	u.Status = "active"
	u.Password = "123"
	u.Password = encrypt([]byte(config.GetString("CRYPT_KEY")), u.Password)
	u.IsDefaultPassword = true
	return
}

// encrypt string to base64 crypto using AES
func encrypt(key []byte, text string) string {
	// key := []byte(keyText)
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}