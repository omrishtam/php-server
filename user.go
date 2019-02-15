package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)

type UserHandler struct {}

func (h UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"]
	fmt.Fprintf(w, "%s", userId)
}