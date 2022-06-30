package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Payload struct {
	Email string
	ID primitive.ObjectID
}


type Claims struct {
	Email string 						`json:"email"`
	ID primitive.ObjectID 	`json:"id"`
	jwt.StandardClaims
}

var JWT_SECRET string


func GenerateJwtToken(payload Payload) (string, error) {

	if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == ""  {
		log.Fatal("JWT_SECRET WAS NOT provided")
	}

	key := []byte(JWT_SECRET)

	expirationTime := time.Now().Add(7 * 24 * 60 * time.Minute)

	claims := &Claims {
		ID: payload.ID,
		Email: payload.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := unsignedToken.SignedString(key)

	if err != nil {
		return "", err
	}
	return signedToken, nil

}


func VerifyToken(strToken string) (*Claims, error) {

	if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
		log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	}

	key := []byte(JWT_SECRET)

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(strToken, claims, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, fmt.Errorf("invalid token signature")
		}
	}

	if !token.Valid {
		return claims, fmt.Errorf("invalid token")
	}

	return claims, nil
}