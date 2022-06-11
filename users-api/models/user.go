package models

import (
	"encoding/json"
	"io"
	"github.com/go-playground/validator/v10"
)

type User struct {
	ID uint64 `json:"id"`
	Name string `json:"name" validate:"required,min=6"`
	Email string `json:"email" validate:"required"`
}

type Models struct {
	Users []*User
}


func (u *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(u)
}


func(u *User) Validate() error{
	validate := validator.New()
	return validate.Struct(u)
}
