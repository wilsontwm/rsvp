package controllers

import (
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var viewPath = "views"
var templates *template.Template
var restURL *url.URL
var apiURL string
var appURL string
var appName string
var appVersion string
var store *sessions.CookieStore
var cookieHashKey []byte
var cookieBlockKey []byte
var sCookie *securecookie.SecureCookie

func init() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Print("Error loading .env file", err)
	}

	templates, _ = GetTemplates()
	appName = os.Getenv("app_name")
	appURL = os.Getenv("app_url")
	apiURL = os.Getenv("app_api_url")
	appVersion = os.Getenv("app_version")
	restURL, _ = url.ParseRequestURI(apiURL)
	store = sessions.NewCookieStore([]byte(os.Getenv("session_key")))
	sCookie = securecookie.New(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))
}

func GetTemplates() (templates *template.Template, err error) {
	var allFiles []string

	funcMap := template.FuncMap{
		"NL2BR": func(value string) string {
			text := template.HTMLEscapeString(value)
			return strings.Replace(text, "\n", "<br>", -1)
		},
		"safeHTML": func(value string) template.HTML {
			return template.HTML(value)
		},
		"truncate": func(value string, limit int) string {
			runes := []rune(value)
			if len(runes) > limit {
				return string(runes[:limit]) + "..."
			}
			return value
		},
	}
	// Loop through all the files in the views folder including subfolders
	err = filepath.Walk(viewPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			allFiles = append(allFiles, path)
		}

		return nil
	})

	if err != nil {
		log.Print("Error walking the file path", err)
	}

	templates, err = template.New("").Funcs(funcMap).ParseFiles(allFiles...)

	if err != nil {
		log.Print("Error parsing template files", err)
	}

	return
}

func SetCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string, cookieValue string) {
	cookie := &http.Cookie{
		Name:  cookieName,
		Value: cookieValue,
		Path:  "/",
		// true means no scripts, http requests only
		HttpOnly: false,
	}

	http.SetCookie(w, cookie)
}

func SetEncodedCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string, cookieValue string) {
	value := cookieValue

	if encoded, err := sCookie.Encode(cookieName, value); err == nil {
		cookie := &http.Cookie{
			Name:  cookieName,
			Value: encoded,
			Path:  "/",
			// true means no scripts, http requests only
			HttpOnly: false,
		}

		http.SetCookie(w, cookie)
	}
}

func ClearCookieHandler(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func ReadCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string) (cookieValue string) {
	cookie, err := r.Cookie(cookieName)

	if err == nil {
		cookieValue = cookie.Value
	}

	return
}

func ReadEncodedCookieHandler(w http.ResponseWriter, r *http.Request, cookieName string) (cookieValue string) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		if err = sCookie.Decode(cookieName, cookie.Value, &cookieValue); err == nil {
			return
		}
	}

	return
}

func CheckAuthenticatedRequest(w http.ResponseWriter, r *http.Request, responseCode int) bool {
	if responseCode == http.StatusUnauthorized {
		ClearCookieHandler(w, "auth")
		ClearCookieHandler(w, "name")
		return false
	}

	return true
}
