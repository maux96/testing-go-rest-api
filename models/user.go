package models

type User struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}
