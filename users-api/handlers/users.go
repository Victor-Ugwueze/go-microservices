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

	"github.com/Victor-Ugwueze/go-microservices/users-api/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Hello is a simple handler
type UserModel struct {
	l *log.Logger
	C *mongo.Collection;
}


type Input struct {
	ID      string `bson:"_id,omitempty"`
	Name     string   `bson:"name,omitempty"`
	Email    string          `bson:"email,omitempty"`
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
			http.Error(rw, "An Error " ,http.StatusServiceUnavailable)
			return
		}

		res, err := json.Marshal(allUsers)


		if err != nil {
			http.Error(rw, "An Error " ,http.StatusServiceUnavailable)
			return
		}
		// write the response
		rw.Write(res)
}

func addUser(user *models.User, um *UserModel) (primitive.ObjectID ,error) {
	result, err := um.C.InsertOne(context.TODO(), user)
	return result.InsertedID.(primitive.ObjectID), err
}



func create(rw http.ResponseWriter, r *http.Request, h *UserModel)  {
	rw.Header().Set("Content-Type", "application/json")
	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		http.Error(rw, "Server error" ,http.StatusServiceUnavailable)
		return
	}

	var user models.User

	fields := bson.D{
		{Key: "_id", Value: 1},
		{Key: "email", Value: 1},
		{Key: "name", Value: 1 },
	}


	opts := options.FindOne()
	opts.SetProjection(fields)


  err = h.C.FindOne(context.TODO(), bson.M{"email": u.Email }, opts).Decode(&user)

	fmt.Println(err, user)

	if err != nil {
		http.Error(rw, "Bad request" ,http.StatusServiceUnavailable)
		return
	}

	 id, err := addUser(&u, h)
	
	if err != nil {
		http.Error(rw, "Bad request" ,http.StatusServiceUnavailable)
		return
	}

	result := models.User{ID: id, Name: u.Name, Email: u.Email }

	err = json.NewEncoder(rw).Encode(result)

	if err != nil {
		http.Error(rw, "Unable to handle request" ,http.StatusServiceUnavailable)
	}
}

func delete(id int, r *http.Request, rw http.ResponseWriter, h *UserModel) {

	// index, err := findProduct(h.Models, id)

	// if err != nil  {
	// 	http.NotFound(rw, r)
	// 	return
	// }
	// h.Models.Users = append(h.Models.Users[:index], h.Models.Users[index+1:]...)
}

func findProduct(h UserModel, id int)(int, error) {
	// for index, u := range models.Users {
	// 	if id == int(u.ID) {
	// 		return index, nil
	// 	}
	// }
	return -1, fmt.Errorf("User not found")
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






// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *UserModel) CreateUsers(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("requests", r.URL.Path)
	create(rw, r, h)
}


// ServeHTTP implements the go http.Handler interface
// https://golang.org/pkg/net/http/#Handler
func (h *UserModel) DeleteUsers(rw http.ResponseWriter, r *http.Request) {
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


	func(h *UserModel) ValidateUserData(next http.Handler) http.Handler {
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