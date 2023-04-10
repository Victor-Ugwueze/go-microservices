package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	models "github.com/Victor-Ugwueze/go-microservices/account-service/models"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginRequest struct {
	Email     string     `json:"email"`
	Password  string     `json:"password"`
}

type LoginResponse struct {
	Token     string     `json:"token"`
}

type AuthHandler struct {
	db *mongo.Collection
}


func NewAuthHandler(db *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}


func GenerateJwtToken(C *AuthHandler) (rw http.ResponseWriter, r *http.Request) {

	userPayload := &LoginRequest{}

	err := json.NewDecoder(r.Body).Decode(&userPayload)

	if err != nil {
		http.Error(rw, "An Error occurred", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	filter := bson.M{"email": userPayload.Email}
	var user  models.User
	err = C.db.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		http.Error(rw, "Unauthorized", http.StatusNotFound)
		return
	}

		
	// Create the Claims

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		Issuer:    "auth-svc",
		Subject:  user.ID.Hex(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		http.Error(rw, "An Error occurred", http.StatusInternalServerError)
		return
	}

	 json.NewEncoder(rw).Encode(&LoginResponse{Token: tokenString})	
	 return
}


func ValidateToken(token string) (bool, error) {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

