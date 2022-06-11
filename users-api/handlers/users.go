
// Package classification of Users API
//
// Documentation for Users API
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
// Consumes:
// - application/json
//
// Produces:
// - application/json
// swagger:meta
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"vic/models"

	"github.com/gorilla/mux"
)

// Hello is a simple handler
type Hello struct {
	l *log.Logger
	Models *models.Models
}


// NewHello creates a new hello handler with the given logger
func Newusers(l *log.Logger, models *models.Models) (*Hello) {
	return &Hello{l, models}
}

// NewHello creates a new hello handler with the given logger
func list(rw http.ResponseWriter, h *Hello) {
		rw.Header().Set("Content-Type", "application/json")
		res, _ := json.Marshal(h.Models.Users)
		// write the response
		rw.Write(res)
}

func addUser(u *models.User, h *Hello, w io.Reader) error {
	id := len(h.Models.Users)
  u.ID = uint64(id + 1)
	h.Models.Users = append(h.Models.Users, u)

	return nil
}



func create(rw http.ResponseWriter, r *http.Request, h *Hello)  {
	val := r.Context().Value(KeyUser{}).(models.User)
	var err = addUser(&val, h, r.Body)
	rw.Header().Set("Content-Type", "application/json")
	

	if err != nil {
		http.Error(rw, "Bad request" ,http.StatusServiceUnavailable)
		return
	}
	err = val.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusServiceUnavailable)
	}
}

func delete(id int, r *http.Request, rw http.ResponseWriter, h *Hello) {

	index, err := findProduct(h.Models, id)

	if err != nil  {
		http.NotFound(rw, r)
		return
	}
	h.Models.Users = append(h.Models.Users[:index], h.Models.Users[index+1:]...)
}

func findProduct(models *models.Models, id int)(int, error) {
	for index, u := range models.Users {
		if id == int(u.ID) {
			return index, nil
		}
	}
	return -1, fmt.Errorf("User not found")
}

func update(id int, r *http.Request, rw http.ResponseWriter, h *Hello) {
	rw.Header().Set("Content-Type", "application/json")
	val := r.Context().Value(KeyUser{}).(models.User)

  index, err := findProduct(h.Models, id)

	if err != nil  {
		http.NotFound(rw, r)
		return
	}

	val.ID = uint64(id)
	h.Models.Users[index] = &val
	val.ToJSON(rw)

}

// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *Hello) ListUsers(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	list(rw, h)
}


// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *Hello) UpdateUsers(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to pass id", http.StatusBadRequest)
		return
	}
	update(id, r, rw, h)
}






// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *Hello) CreateUsers(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	create(rw, r, h)
}


// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *Hello) DeleteUsers(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Unable to pass id", http.StatusBadRequest)
		return
	}
	
		delete(id, r, rw, h)
	}


	type KeyUser struct {}


	func(h *Hello) ValidateUserData(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			user := &models.User{}
			e := json.NewDecoder(r.Body)
			var err = e.Decode(user)

			if err != nil {
				http.Error(rw, "Unable to handle request" ,http.StatusServiceUnavailable)
				return
			}

			err = user.Validate()

			if err != nil {
				h.l.Println("Error val", err)
				http.Error(rw, fmt.Sprintf("Error validating user: %s", err) ,http.StatusServiceUnavailable)
				return
			}

			ctx := context.WithValue(r.Context(), KeyUser{}, user)
			r = r.WithContext(ctx)
			next.ServeHTTP(rw, r)
		})
	}