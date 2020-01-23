package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"rsvp/utils"
)

// Delete the attendee
var DeleteAttendeeSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	vars := mux.Vars(r)

	// Set the URL path
	restURL.Path = "/api/auth/attendees/" + vars["id"] + "/delete"
	urlStr := restURL.String()

	session, err := utils.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Set the input data
	jsonData := map[string]interface{}{}

	response, err := utils.SendAuthenticatedRequest(urlStr, auth, jsonData, "DELETE")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(data, &resp)
		utils.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}
