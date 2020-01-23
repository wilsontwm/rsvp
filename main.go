package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"rsvp/api"
	"rsvp/controllers"
	"rsvp/middleware"
)

func main() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	router := mux.NewRouter()

	// API routes
	apiRoutes := router.PathPrefix("/api").Subrouter()

	// General routes
	apiRoutes.HandleFunc("/login", api.Login).Methods("POST")
	apiRoutes.HandleFunc("/signup", api.Signup).Methods("POST")

	// Attendees routes
	apiAttendeesRoutes := apiRoutes.PathPrefix("/attendees").Subrouter()
	apiAttendeesRoutes.HandleFunc("/create", api.CreateAttendees).Methods("POST")

	/*********************************
	 *      Authenticated routes     *
	 *********************************/
	apiAuthenticatedRoutes := apiRoutes.PathPrefix("/auth").Subrouter()
	apiAuthenticatedRoutes.Use(middleware.JwtAuthentication())

	// Attendees routes
	apiAuthAttendeesRoutes := apiAuthenticatedRoutes.PathPrefix("/attendees").Subrouter()
	apiAuthAttendeesRoutes.HandleFunc("/", api.IndexAttendees).Methods("GET")
	apiAuthAttendeesRoutes.HandleFunc("/{id}/delete", api.DeleteAttendee).Methods("DELETE")

	// User routes
	apiUserRoutes := apiAuthenticatedRoutes.PathPrefix("/profile").Subrouter()
	apiUserRoutes.HandleFunc("/get", api.GetProfile).Methods("GET")
	apiUserRoutes.HandleFunc("/edit", api.EditProfile).Methods("PATCH")
	apiUserRoutes.HandleFunc("/edit/password", api.EditPassword).Methods("PATCH")

	// ******************************************************************************* //
	// Web routes
	routes := router.PathPrefix("").Subrouter()

	// General routes
	routes.HandleFunc("/", controllers.HomePage).Methods("GET")
	routes.HandleFunc("/login", controllers.LoginPage).Methods("GET")
	routes.HandleFunc("/login", controllers.LoginSubmit).Methods("POST")
	routes.HandleFunc("/signup", controllers.SignupPage).Methods("GET")
	routes.HandleFunc("/signup", controllers.SignupSubmit).Methods("POST")
	routes.HandleFunc("/logout", controllers.LogoutSubmit).Methods("POST")
	routes.HandleFunc("/rsvp", controllers.RsvpSubmit).Methods("POST")

	// Authenticated routes
	authenticatedRoutes := routes.PathPrefix("/dashboard").Subrouter()
	authenticatedRoutes.HandleFunc("", controllers.DashboardPage).Methods("GET")
	authenticatedRoutes.HandleFunc("/profile/edit", controllers.EditProfilePage).Methods("GET")
	authenticatedRoutes.HandleFunc("/profile/edit", controllers.EditProfileSubmit).Methods("POST")
	authenticatedRoutes.HandleFunc("/profile/edit/password", controllers.EditPasswordSubmit).Methods("POST")
	authenticatedRoutes.HandleFunc("/attendees/{id}/delete", controllers.DeleteAttendeeSubmit).Methods("POST")

	// Asset files
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))
	router.PathPrefix("/storage/").Handler(http.StripPrefix("/storage/", http.FileServer(http.Dir("./storage/"))))

	port := os.Getenv("port")
	if port == "" {
		port = "8000"
	}

	log.Println("Server started and running at port", port)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
