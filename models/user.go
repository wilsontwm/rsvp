package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"rsvp/utils"
	"time"
)

type Token struct {
	UserId uuid.UUID
	Expiry time.Time
	jwt.StandardClaims
}

type User struct {
	Base
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	PassCode string `gorm:"-"`
	Token    string `gorm:"-"`
}

// Login the user
func (user *User) Login(email string, password string) map[string]interface{} {
	var resp map[string]interface{}

	// Get the user by email
	db := GetDB()
	db.Table("users").Where("email = ?", email).First(&user)

	defer db.Close()

	if user.Email == "" {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Invalid email address or password.")
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		// If password does not match
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
			resp = utils.Message(false, http.StatusUnprocessableEntity, "Invalid email address or password.")
		} else {
			// Password matches
			user.Password = "" // remove the password

			// Create new JWT token for the newly registered account
			expiry := time.Now().Add(time.Hour * 2) // Only valid for 2 hours
			tk := &Token{UserId: user.ID, Expiry: expiry}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
			tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
			user.Token = tokenString

			resp = utils.Message(true, http.StatusOK, "You have successfully logged in.")
			resp["data"] = user
		}
	}

	return resp
}

// Validate the incoming details for signup
func (user *User) ValidateSignup() (map[string]interface{}, bool) {
	var resp map[string]interface{}

	// Email must be unique
	temp := &User{}

	// Check for errors and duplicate emails
	db := GetDB()
	err := db.Table("users").Where("email = ?", user.Email).First(temp).Error

	defer db.Close()

	if err != nil && err != gorm.ErrRecordNotFound {
		resp = utils.Message(false, http.StatusInternalServerError, "Connection error. Please retry.")
		return resp, false
	}

	if temp.Email != "" {
		resp = utils.Message(false, http.StatusUnprocessableEntity, "Email address has already been taken.")
		return resp, false
	}

	resp = utils.Message(true, http.StatusOK, "Input has been validated.")
	return resp, true
}

// Register a new user
func (user *User) Create() map[string]interface{} {
	var resp map[string]interface{}

	// Validate the account first
	if resp, ok := user.ValidateSignup(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	db := GetDB()

	defer db.Close()
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, "Failed to create account: "+err.Error())
		return resp
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		resp = utils.Message(false, http.StatusInternalServerError, "Failed to create account: "+err.Error())
		return resp
	}

	err := tx.Commit().Error

	if user.ID == uuid.Nil && err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, "Failed to create account: "+err.Error())
		return resp
	}

	user.Password = "" // delete the password

	resp = utils.Message(true, http.StatusOK, "You have successfully signed up.")
	resp["data"] = user

	return resp
}

func (user *User) Edit() map[string]interface{} {
	db := GetDB()

	defer db.Close()

	db.Model(&user).Update(map[string]interface{}{
		"Name":  user.Name,
		"Email": user.Email,
	})

	resp := utils.Message(true, http.StatusOK, "Successfully updated profile.")
	resp["data"] = user

	return resp
}

func (user *User) EditPassword() map[string]interface{} {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	password := string(hashedPassword)

	db := GetDB()

	defer db.Close()

	db.Model(&user).Update(map[string]interface{}{
		"Password": password,
	})

	resp := utils.Message(true, http.StatusOK, "Successfully updated password.")

	return resp
}

func getUser(user *User) *User {
	if user.Email == "" {
		return nil
	}

	user.Password = ""

	return user
}

// Get the user by ID
func GetUser(u uuid.UUID) *User {
	user := &User{}
	db := GetDB()

	defer db.Close()

	db.Table("users").
		Where("id = ?", u).
		First(user)

	return getUser(user)
}
