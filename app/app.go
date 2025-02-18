package app

import (
	"Go_Assignment/auth"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func HandleRoutes() {
	router := mux.NewRouter()

	// Authentication routes
	router.HandleFunc("/login", Login).Methods(http.MethodGet)

	// End point to refresh authentication JWT token
	router.HandleFunc("/refresh", RefreshToken).Methods(http.MethodGet)

	// Define all routes and endpoints
	// End point to fetch All users
	router.HandleFunc("/users", auth.AuthMiddleware(FetchAllUsers)).Methods(http.MethodGet)

	// End point to fetch specific user
	router.HandleFunc("/user", auth.AuthMiddleware(FetchSpecificUser)).Methods(http.MethodGet)

	// End point to Save user
	router.HandleFunc("/saveUser", SaveUser).Methods(http.MethodPost)

	// End point to Update user
	router.HandleFunc("/updateUser", UpdateUser).Methods(http.MethodPut)

	// End point to delete a specific user
	router.HandleFunc("/deleteUser", DeleteUser).Methods(http.MethodDelete)

	// Start the server at 8080 port
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
