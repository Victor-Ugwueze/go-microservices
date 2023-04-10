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
	"log"
	"net/http"
	"strconv"

	"github.com/Victor-Ugwueze/go-microservices/users-api/helpers"
	"github.com/Victor-Ugwueze/go-microservices/users-api/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Hello is a simple handler
type UserModel struct {
	l *log.Logger
	C *mongo.Collection;
}

type KeyUser struct {}



type AuthResponse struct {
	User models.User `json:"user"`
	Token string     `json:"token"`
}


type LoginRequest struct {
	Email     string     `json:"email"`
	Password  string     `json:"password"`
}

// NewHello creates a new hello handler with the given logger
func Newusers(l *log.Logger, C *mongo.Collection) (*UserModel) {
	return &UserModel{l, C}
}

// NewHello creates a new hello handler with the given logger
func list(rw http.ResponseWriter, um *UserModel) {
		rw.Header().Set("Content-Type", "application/json")
		ctx := context.TODO()
		allUsers := []models.User{}

		userCursor, err := um.C.Find(ctx, bson.M{})

		if err != nil {
			http.Error(rw, "An Error occurred" ,http.StatusServiceUnavailable)
			fmt.Println(err)
			return
		}

		err = userCursor.All(ctx, &allUsers)

				if err != nil {
			http.Error(rw, "An Error ocurred" ,http.StatusServiceUnavailable)
			return
		}

		res, err := json.Marshal(allUsers)

		if err != nil {
			http.Error(rw, "An Error " ,http.StatusServiceUnavailable)
			return
		}
		rw.Write(res)
}



func findUserByEmail(email string, h *UserModel) (*models.User, error) {
	result := h.C.FindOne(context.TODO(), bson.M{"email": email })

	var u = &models.User{}

	err := result.Decode(u)

	if err != nil {
		return u, err
	}
	return u, nil
}


func create(u *models.User, rw http.ResponseWriter, h *UserModel) (error)  {
	var err = u.HashPassword()
	result, err := h.C.InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	u.ID = result.InsertedID.(primitive.ObjectID)
	return  nil
}


func checkEmailExists(u *models.User, h *UserModel) (bool, error) {
	count, err := h.C.CountDocuments(context.TODO(), bson.M{"email": u.Email })

	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, nil
	}
}


// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *UserModel) Signup(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)

	user := r.Context().Value(KeyUser{}).(*models.User)

	userFound, err := checkEmailExists(user, h)

	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusInternalServerError)
		return
	}

	if userFound {
		http.Error(rw, "Email already exists", http.StatusConflict)
		return
	}

	 err = create(user, rw, h)

	 if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusInternalServerError)
		return
	}

	token, err := helpers.GenerateJwtToken(helpers.Payload{ ID: user.ID, Email: user.Email })

	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusInternalServerError)
		return
	}


  user.Password = ""
	json.NewEncoder(rw).Encode(AuthResponse{ User: *user, Token: token})

	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusInternalServerError)
		return
	}
}


// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *UserModel) Login(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	rw.Header().Set("Content-Type", "application/json")

	var payload  = LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(&payload)


	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusUnauthorized)
		return
	}

	user, err := findUserByEmail(payload.Email, h)

	if err != nil {
		http.Error(rw, "Authentication failed" ,http.StatusUnauthorized)
		return
	}

	err = user.ComparePassword(payload.Password)

	if err != nil {
		http.Error(rw, "Authentication failed" ,http.StatusUnauthorized)
		return
	}


	token, err := helpers.GenerateJwtToken(helpers.Payload{ID: user.ID, Email: user.Email })

	 if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusInternalServerError)
		return
	}


	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusInternalServerError)
		return
	}

	user.Password = ""
	json.NewEncoder(rw).Encode(AuthResponse{ User: *user, Token: token})

	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusInternalServerError)
		return
	}
}




func update(id int, r *http.Request, rw http.ResponseWriter, h *UserModel) {
	rw.Header().Set("Content-Type", "application/json")
	// val := r.Context().Value(KeyUser{}).(models.User)

  // index, err := findProduct(h.Models, id)

	// if err != nil  {
	// 	http.NotFound(rw, r)
	// 	return
	// }

	// val.ID = uint64(id)
	// h.Models.Users[index] = &val
	// val.ToJSON(rw)

}

// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *UserModel) ListUsers(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	list(rw, h)
}


// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *UserModel) UpdateUsers(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to pass id", http.StatusBadRequest)
		return
	}
	update(id, r, rw, h)
}






	func(h *UserModel) ValidateUserData(next http.Handler) http.Handler {

		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Content-Type", "application/json")
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
	