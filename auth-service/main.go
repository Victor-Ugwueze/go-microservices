package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/victor-ugwueze/go-microservices/auth-service/handlers"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {


	l := log.New(os.Stdout, "users-service-api ", log.LstdFlags)

  err := godotenv.Load()

  if err != nil {
    l.Fatal("Error loading .env file")
  }

	mongoURI := flag.String("mongoURI", "mongodb://root:password@localhost:27017", "Database hostname url")
	serverAddr := flag.String("serverAddr", "localhost", "Network address")
	serverPort := flag.Int("serverPort", 9000, "Port")

	co := options.Client().ApplyURI(*mongoURI)



	// establish database connection
	client, err := mongo.NewClient(co)

	// create handler
	authHandler := handlers.NewAuth(client)

	if err != nil {
		l.Fatal(err)
	}


	serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)

	http.HandleFunc("/generate-token", authHandler.GenerateToken)
	http.HandleFunc("/validate-token", handler.validateToken)

	s := http.Server{
		Addr:  serverURI,
	}

	err = s.ListenAndServe()
	if err != nil {
		 l.Printf("There was an error listening on port :9000", err)
	}

}
