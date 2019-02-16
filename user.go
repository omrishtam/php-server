package main

import (
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"net/http"
	"encoding/json"
	"context"
)

// User is a structure used for serializing/deserialzing user data
type User struct {
	ID string `json:"_id"`
	Name string `json:"name"`
	Admin bool `json:"admin"`
}

// UserHandler is a structure used for handling requests for user related actions
type UserHandler struct {
	Collection *mongo.Collection
}

// GetUserHandler Gets a GET request with a user's id and responds with the
// user's data from the database
func (h UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	filter := bson.D{{"_id", userID}}
	var user User

	err := h.Collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseUser, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Write(responseUser)
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
