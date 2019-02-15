package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
	"encoding/json"
)

// User is a structure used for serializing/deserialzing user data
type User struct {
	ID string `json:"id"`
}

// UserHandler is a structure used for handling requests for user related actions
type UserHandler struct {}

// GetUserHandler Gets a GET request with a user's id and responds with the
// user's data from the database
func (h UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := ""
	userID = vars["id"]

	fmt.Fprintf(w, "%s", userID)
}

// AddUserHandler Gets a POST request with user's data in the request's body and saves
// it in the database and responds with the new user's data
func (h UserHandler) AddUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var user User
	err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UpdateUserHandler Gets a PUT request with user's data in the request's body
// and updates it in the database and responds with the updated user's data
func (h UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	if userID == "" {

	}

	decoder := json.NewDecoder(r.Body)

	var user User
	err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// DeleteUserHandler Gets a DELETE request with a user's id and responds with 
// true or false wether the user was deleted from the database or not
func (h UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	userID = ""

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(userID))
}