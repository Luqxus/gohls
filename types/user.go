package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID
	UID       string
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

func (u *User) FormatResponse() *ResponseUser {
	return &ResponseUser{
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

// TODO: verify input fields
type UserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// TODO: verify input fields
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ResponseUser struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
