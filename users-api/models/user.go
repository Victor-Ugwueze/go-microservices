package models

import (
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      	primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     	string             `json:"name" validate:"required"`
	Email 		string             `json:"email" validate:"required,email"`
	Password 	string						 `json:"password,omitempty" validate:"required,min=8"`
}

type AuthResponse struct {
	ID      	primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     	string             `json:"name" validate:"required"`
	Email 		string             `json:"email" validate:"required,email"`
}

func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) FromJSON(r *http.Request) error {
	e := json.NewDecoder(r.Body)
	return e.Decode(u)
}


func(u *User) Validate() error{
	validate := validator.New()
	return validate.Struct(u)
}

func(u *User) HashPassword() error{
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func(u *User) ComparePassword(password string) error{
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))

	if err != nil {
		return err
	}
	return nil
}
