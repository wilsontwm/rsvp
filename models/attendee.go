package models

import (
	"net/http"
	"rsvp/utils"
	"strconv"
	"strings"
)

type Attendee struct {
	Base
	Name            string `gorm:"not null"`
	Email           string
	Phone           string
	UpdatedAtString string `gorm:"-"`
}

// Validate the incoming details for signup
func (attendee *Attendee) ValidateAttendees(name string, email string, phone string, names []string, emails []string, phones []string) (map[string]interface{}, bool) {
	var resp map[string]interface{}

	for _, n := range names {
		if strings.TrimSpace(n) == "" {
			resp = utils.Message(false, http.StatusUnprocessableEntity, "Validation error: Cannot leave attendee name empty..")
			return resp, false
		}
	}

	resp = utils.Message(true, http.StatusOK, "Input has been validated.")
	return resp, true
}

// Create multiple attendees
func (attendee *Attendee) CreateMultiple(name string, email string, phone string, names []string, emails []string, phones []string) map[string]interface{} {
	var resp map[string]interface{}

	// Validate the account first
	if resp, ok := attendee.ValidateAttendees(name, email, phone, names, emails, phones); !ok {
		return resp
	}

	if err := CreateAttendeesTransaction(name, email, phone, names, emails, phones); err != nil {
		resp = utils.Message(false, http.StatusInternalServerError, err.Error())
		return resp
	}

	resp = utils.Message(true, http.StatusOK, "Thank you! We are looking forward to seeing you!")

	return resp
}

// The transaction to create attendees in bulk
func CreateAttendeesTransaction(name string, email string, phone string, names []string, emails []string, phones []string) error {
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

	mainAttendee := Attendee{Name: name, Email: email, Phone: phone}
	if err := tx.Create(&mainAttendee).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Loop through all the accompanies and create
	for i := 0; i < len(names) && i < len(emails) && i < len(phones); i++ {
		attendee := Attendee{Name: names[i], Email: emails[i], Phone: phones[i]}
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
		Select("attendees.*, TO_CHAR(attendees.updated_at, '" + utils.DateTimeSQLFormat + "') as updated_at_string").
		Where("deleted_at is NULL").
		Order("updated_at DESC").
		Find(&attendees)

	resp = utils.Message(true, http.StatusOK, "You have successfully retrieved "+strconv.Itoa(len(attendees))+" attendees.")
	resp["data"] = attendees

	return resp
}
