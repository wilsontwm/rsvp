package controllers

import (
	"encoding/json"
	"github.com/gorilla/csrf"
	"io/ioutil"
	"net/http"
	"rsvp/utils"
	"strings"
	"time"
)

// Edit profile page
var EditProfilePage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// Set the URL path
	restURL.Path = "/api/auth/profile/get"
	urlStr := restURL.String()

	session, err := utils.GetSession(store, w, r)

	jsonData := make(map[string]interface{})
	auth := ReadEncodedCookieHandler(w, r, "auth")
	response, err := utils.SendAuthenticatedRequest(urlStr, auth, jsonData, "GET")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		respData, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(respData, &resp)

		if resp["success"].(bool) {

			data := map[string]interface{}{
				"title":          "Edit profile",
				"appName":        appName,
				"appVersion":     appVersion,
				"year":           time.Now().Format("2006"),
				"user":           resp["data"],
				csrf.TemplateTag: csrf.TemplateField(r),
			}

			data, err = InitializePage(w, r, store, data)

			err = templates.ExecuteTemplate(w, "edit_profile_html", data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else {
			utils.SetErrorSuccessFlash(session, w, r, resp)
			// Redirect back to the dashboard page
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		}
	}
}

// Post request for editing the profile for the current user
var EditProfileSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/auth/profile/edit"
	urlStr := restURL.String()

	session, err := utils.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	email := strings.TrimSpace(r.Form.Get("email"))

	// Set the input data
	jsonData := map[string]interface{}{
		"name":  name,
		"email": email,
	}

	response, err := utils.SendAuthenticatedRequest(urlStr, auth, jsonData, "PATCH")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(data, &resp)

		if resp["success"].(bool) {
			// Need to reset the cookie that store name
			userData := resp["data"].(map[string]interface{})
			SetCookieHandler(w, r, "name", userData["Name"].(string))
		}

		utils.SetErrorSuccessFlash(session, w, r, resp)

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// Post request to edit the password for the current user
var EditPasswordSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/auth/profile/edit/password"
	urlStr := restURL.String()

	session, err := utils.GetSession(store, w, r)

	// Get the auth info for edit profile
	auth := ReadEncodedCookieHandler(w, r, "auth")

	// Get the input data from the form
	r.ParseForm()
	password := strings.TrimSpace(r.Form.Get("password"))
	retype_password := strings.TrimSpace(r.Form.Get("retype_password"))

	// Set the input data
	jsonData := map[string]interface{}{
		"password":        password,
		"retype_password": retype_password,
	}

	response, err := utils.SendAuthenticatedRequest(urlStr, auth, jsonData, "PATCH")

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
