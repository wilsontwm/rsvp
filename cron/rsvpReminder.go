package cron

import (
	"fmt"
	"github.com/robfig/cron"
	"rsvp/models"
	"rsvp/utils"
	"time"
)

// Schedule a task on 19 September to send RSVP reminder mail
func ScheduleRSVPReminder() {
	c := cron.New()

	c.AddFunc("* * 19 9 *", func() {
		fmt.Println("[Job] Sending RSVP reminders...")
		SendRSVPReminderMail()
	})

	c.Start()
}

// Send RSVP reminder mail
func SendRSVPReminderMail() {
	year, _, _ := time.Now().Date()
	if year == 2020 {
		// Get all the registered attendees with emails
		attendees := []models.Attendee{}

		db := models.GetDB()
		defer db.Close()

		db.Table("attendees").
			Where("attendees.deleted_at is NULL").
			Where("attendees.email is not NULL").
			Where("attendees.email <> ''").
			Find(&attendees)

		appName := "Wilson & Shu Zhen"
		subject := appName + " - We'll see you at our wedding soon!"

		for _, attendee := range attendees {
			receiver := attendee.Email
			r := utils.NewRequest([]string{receiver}, subject)
			go r.Send("views/mail/reminder.html", map[string]string{"appName": appName, "name": attendee.Name})
		}

	}
}
