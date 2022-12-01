package models

type User struct {
	Id               int    `json:"id,omitempty"`
	Name             string `json:"name,omitempty"`
	Password         string `json:"password,omitempty"`
	Birthdate        string `json:"birthdate,omitempty"`
	ProfilePictureId string `json:"profilePictureId,omitempty"`
}
