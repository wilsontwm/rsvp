package controllers

import (
	"encoding/json"
	"github.com/gorilla/csrf"
	//"github.com/gorilla/mux"
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

	data, err := utils.InitializePage(w, r, store, data)
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

	data, err := utils.InitializePage(w, r, store, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "login_html", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	data, err := utils.InitializePage(w, r, store, data)
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
