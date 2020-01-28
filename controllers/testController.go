package controllers

import (
	"net/http"
	"rsvp/cron"
)

// GET: Test functionality
var TestFunction = func(w http.ResponseWriter, r *http.Request) {
	cron.SendRSVPReminderMail()
}
