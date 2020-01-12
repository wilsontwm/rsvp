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

// Home page
var HomePage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":          "Home",
		"appName":        appName,
		"appVersion":     appVersion,
		"year":           time.Now().Format("2006"),
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	data, err := InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "index_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Login page
var LoginPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":          "Login",
		"appName":        appName,
		"appVersion":     appVersion,
		"year":           time.Now().Format("2006"),
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	data, err := InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "login_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST: Login functionality
var LoginSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// Set the URL path
	restURL.Path = "/api/login"
	urlStr := restURL.String()

	session, err := utils.GetSession(store, w, r)

	// Get the input data from the form
	r.ParseForm()
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))

	// Set the input data
	jsonData := map[string]interface{}{
		"email":    email,
		"password": password,
	}

	url := r.Header.Get("Referer")
	response, err := utils.SendRequest(urlStr, jsonData, "POST")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(data, &resp)

		// If login is authenticated
		if resp["success"].(bool) {
			userData := resp["data"].(map[string]interface{})
			// Store the user token in the cookie
			SetEncodedCookieHandler(w, r, "auth", userData["Token"].(string))
			SetCookieHandler(w, r, "name", userData["Name"].(string))
			SetCookieHandler(w, r, "id", userData["ID"].(string))
			url = "/dashboard"
		} else {
			utils.SetErrorSuccessFlash(session, w, r, resp)
		}

		// Redirect back to the previous page
		http.Redirect(w, r, url, http.StatusFound)
	}
}

// Signup page
var SignupPage = func(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"title":          "Signup",
		"appName":        appName,
		"appVersion":     appVersion,
		"year":           time.Now().Format("2006"),
		csrf.TemplateTag: csrf.TemplateField(r),
	}

	data, err := InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "signup_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// POST: Signup functionality
var SignupSubmit = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}

	// Set the URL path
	restURL.Path = "/api/signup"
	urlStr := restURL.String()

	session, err := utils.GetSession(store, w, r)

	// Get the input data from the form
	r.ParseForm()
	name := strings.TrimSpace(r.Form.Get("name"))
	email := strings.TrimSpace(r.Form.Get("email"))
	password := strings.TrimSpace(r.Form.Get("password"))
	retype_password := strings.TrimSpace(r.Form.Get("retype_password"))
	passcode := strings.TrimSpace(r.Form.Get("passcode"))

	// Set the input data
	jsonData := map[string]interface{}{
		"email":           email,
		"password":        password,
		"retype_password": retype_password,
		"name":            name,
		"passcode":        passcode,
	}

	response, err := utils.SendRequest(urlStr, jsonData, "POST")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		// Parse it to json data
		json.Unmarshal(data, &resp)

		utils.SetErrorSuccessFlash(session, w, r, resp)

		if resp["success"].(bool) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Redirect back to the previous page
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	}
}

// POST: Logout functionality
var LogoutSubmit = func(w http.ResponseWriter, r *http.Request) {
	session, _ := utils.GetSession(store, w, r)

	ClearCookieHandler(w, "auth")
	ClearCookieHandler(w, "name")

	session.AddFlash("You have successfully logged out.", "success")
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusFound)
}

// Dashboard page
var DashboardPage = func(w http.ResponseWriter, r *http.Request) {
	var resp map[string]interface{}
	// Set the URL path
	restURL.Path = "/api/auth/attendees/"
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
				"title":          "Dashboard",
				"appName":        appName,
				"appVersion":     appVersion,
				"year":           time.Now().Format("2006"),
				"attendees":      resp["data"],
				csrf.TemplateTag: csrf.TemplateField(r),
			}

			data, err = InitializePage(w, r, store, data)

			err = templates.ExecuteTemplate(w, "dashboard_html", data)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

		} else {
			utils.SetErrorSuccessFlash(session, w, r, resp)
			// Redirect back to the home page
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

}
