package controllers

import "net/http"

//
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Creating a user"))
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting all the users"))
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Getting user"))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Deleting user"))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Updating user"))
}
