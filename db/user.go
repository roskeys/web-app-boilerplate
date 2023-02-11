package db

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/lithammer/shortuuid/v4"
	"github.com/roskeys/app/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID      string `gorm:"column:uid"`
	Email    string `gorm:"column:email"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (User) TableName() string {
	return "user"
}

var usernameRegex = regexp.MustCompile("[a-zA-Z0-9 ]{3,64}")

func CreateNewUser(email, username, password string) (string, utils.ErrorResponse) {
	email = strings.Trim(email, " ")
	username = strings.Trim(username, " ")
	password = strings.Trim(password, " ")
	if _, err := mail.ParseAddress(email); err != nil {
		return "", utils.INVALID_EMAIL_ADDRESS
	}
	existingUser := GetUserByEmail(email)
	if len(existingUser.UID) > 1 {
		return "", utils.USER_EXISTS_FOR_SIGNUP
	}
	if !usernameRegex.MatchString(username) || len(username) > 64 {
		return "", utils.INVALID_USERNAME
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", utils.FAILED_TO_CREATE_USER
	}
	user := User{
		Email:    email,
		Username: username,
		Password: string(passwordHash),
		UID:      shortuuid.New()[:8],
	}
	result := DB.Create(&user)
	if result.Error != nil {
		return "", utils.FAILED_TO_CREATE_USER
	}
	return user.UID, utils.ErrorResponse{}
}

func GetUserByEmail(email string) User {
	var user User
	if result := DB.Where("email = ?", email).First(&user); result.Error != nil {
		return user
	}
	return user
}
