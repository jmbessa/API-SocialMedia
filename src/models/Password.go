package models

//represents the format of the password update request
type Password struct {
	New     string `json:"new"`
	Current string `json:"current"`
}
