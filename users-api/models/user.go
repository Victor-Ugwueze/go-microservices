package models

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID      	primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     	string             `json:"name,omitempty"`
	Email 		string             `json:"email,omitempty"`
}

func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}

func (u *User) FromJSON(r *http.Request) error {
	e := json.NewDecoder(r.Body)
	return e.Decode(&u)
}


func(u *User) Validate() error{
	validate := validator.New()
	return validate.Struct(u)
}
