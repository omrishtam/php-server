package main

import (
	"github.com/gorilla/mux"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"net/http"
	"encoding/json"
	"context"
)

// User is a structure used for serializing/deserialzing user data
type User struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Admin bool `json:"admin" bson:"admin"`
}

// UserHandler is a structure used for handling requests for user related actions
type UserHandler struct {
	Collection *mongo.Collection
}

// GetUsersHandler Gets a GET request with and responds with the
// list of all users' data from the database
func (h UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	findOptions := options.Find()
	filter := bson.D{}
	var users []*User

	cur, err := h.Collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		users = append(users, &user)
	}

	if err := cur.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = cur.Close(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseUsers, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Write(responseUsers)
}

// GetUserHandler Gets a GET request with a user's id and responds with the
// user's data from the database
func (h UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	userOid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	filter := bson.D{bson.E{ Key: "_id", Value: userOid}}
	var user User

	err = h.Collection.FindOne(context.Background(), filter).Decode(&user)
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
	
	user.ID = primitive.NewObjectID()

	insertResult, err := h.Collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	insertResponse, err := json.Marshal(insertResult.InsertedID)

	w.Header().Set("Content-Type", "application/json")
	w.Write(insertResponse)
}

// UpdateUserHandler Gets a PUT request with user's data in the request's body
// and updates it in the database and responds with the updated user's data
func (h UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	userOid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	
	var user User
	err = decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.D{bson.E{ Key: "_id", Value: userOid}}
	update := bson.D{
		{ Key: "$set", Value: bson.D{{ Key: "name", Value: user.Name}}},
		{ Key: "$set", Value: bson.D{{ Key: "admin", Value: user.Admin}}},
	}

	updateResult, err := h.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateResponse, err := json.Marshal(updateResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(updateResponse)
}

// DeleteUserHandler Gets a DELETE request with a user's id and responds with 
// true or false wether the user was deleted from the database or not
func (h UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	userOid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filter := bson.D{bson.E{ Key: "_id", Value: userOid}}
	deleteResult, err := h.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deleteResponse, err := json.Marshal(deleteResult)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(deleteResponse)
}
