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
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Admin bool `json:"admin,omitempty" bson:"admin,omitempty"`
}

// UserHandler is a structure used for handling requests for user related actions
type UserHandler struct {
	UserService
}

// UserService is a structure used for CRUD operations on the users collection 
// in the database
type UserService struct {
	Collection *mongo.Collection
}

// GetAll returns a list of all users in the database that fit the filter
func (s UserService) GetAll(ctx context.Context, filter interface{}, opts ...*options.FindOptions) ([]User, error) {
	var users []User
	cur, err := s.Collection.Find(ctx, filter, opts...)
	defer cur.Close(context.Background())
	if err != nil {
		return nil, err
	}

	for cur.Next(context.Background()) {
		var user User
		err := cur.Decode(&user)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	
	return users, nil
}

// GetOne returns the user in the database that fit the filter
func (s UserService) GetOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (User, error) {
	var user User
	err := s.Collection.FindOne(ctx, filter, opts...).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// InsertOne inserts a user into the users collection in the database and returns its ObjectID
func (s UserService) InsertOne(ctx context.Context, user interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	insertResult, err := s.Collection.InsertOne(ctx, user, opts...)
	if err != nil {
		return nil, err
	}

	return insertResult.InsertedID, nil
}

// UpdateOne updates an existing user in the users collection and returns the updated user's data
func (s UserService) UpdateOne(ctx context.Context,
	filter interface{},
	update interface{},
	opts ...*options.FindOneAndUpdateOptions) (User, error) {
	updateResult := s.Collection.FindOneAndUpdate(ctx, filter, update, opts...)
	
	var updatedUser User
	if updateResult.Err() != nil {
		return User{}, updateResult.Err()
	}

	err := updateResult.Decode(&updatedUser)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

// DeleteOne deletes an existing user in the users collection and returns the deleted user's data
func (s UserService) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) (User, error) {
	deleteResult := s.Collection.FindOneAndDelete(ctx, filter, opts...)
	if deleteResult.Err() != nil {
		return User{}, deleteResult.Err()
	}

	var deletedUser User
	err := deleteResult.Decode(&deletedUser)
	if err != nil {
		return User{}, err
	}

	return deletedUser, nil
}

// GetUsersHandler Gets a GET request and responds with the
// list of all users' data from the database
func (h UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.GetAll(context.Background(), bson.D{})
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

	user, err := h.GetOne(context.Background(), filter)
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
	
	user.ID = primitive.NilObjectID
	user.Admin = false

	if user.Name != "" {
		userID, err := h.InsertOne(context.Background(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		insertResponse, err := json.Marshal(userID)

		w.Header().Set("Content-Type", "application/json")
		w.Write(insertResponse)
	} else {
		http.Error(w, "user name is required", http.StatusBadRequest)
		return
	}
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
	var update []bson.E
	if user.Name != "" {
		update = append(update, bson.E{ Key: "$set", Value: bson.D{{ Key: "name", Value: user.Name}}})
	}

	// TODO: Add a check that the user who made this request is an admin
	update = append(update, bson.E{ Key: "$set", Value: bson.D{{ Key: "admin", Value: user.Admin}}})
	returnDocument := options.ReturnDocument(options.After)
	opts := &options.FindOneAndUpdateOptions{ReturnDocument: &returnDocument}
	updatedUser, err := h.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updateResponse, err := json.Marshal(updatedUser)
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
	deletedUser, err := h.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	deleteResponse, err := json.Marshal(deletedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(deleteResponse)
}
