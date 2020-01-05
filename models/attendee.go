package models

import (
	"net/http"
	"rsvp/utils"
	"strconv"
)

type Attendee struct {
	Base
	Name  string `json:"name" gorm:"not null"`
	Email string `json:"email"`
}

// Create multiple attendees
func (attendee *Attendee) CreateMultiple(name string, email string, names []string, emails []string) map[string]interface{} {
	var resp map[string]interface{}

	if err := CreateAttendeesTransaction(name, email, names, emails); err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, err.Error())
		return resp
	}

	resp = utils.Message(true, http.StatusOK, "Thank you! We are looking forward to seeing you!")

	return resp
}

// The transaction to create attendees in bulk
func CreateAttendeesTransaction(name string, email string, names []string, emails []string) error {
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
		return err
	}

	mainAttendee := Attendee{Name: name, Email: email}
	if err := tx.Create(&mainAttendee).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Loop through all the accompanies and create
	for i := 0; i < len(names) && i < len(emails); i++ {
		attendee := Attendee{Name: names[i], Email: emails[i]}
		if err := tx.Create(&attendee).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (attendee *Attendee) List() map[string]interface{} {
	var resp map[string]interface{}

	attendees := []Attendee{}

	db := GetDB()
	defer db.Close()

	db.Table("attendees").
		Where("deleted_at is NULL").
		Order("updated_at DESC").
		Find(&attendees)

	resp = utils.Message(true, http.StatusOK, "You have successfully retrieved "+strconv.Itoa(len(attendees))+" attendees.")
	resp["data"] = attendees

	return resp
}
