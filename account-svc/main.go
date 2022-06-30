package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Victor-Ugwueze/go-microservices/users-api/handlers"
	"github.com/Victor-Ugwueze/go-microservices/users-api/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)



func Welcome(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(wr, "Welcome here")
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    *models.User
}

func main() {

	l := log.New(os.Stdout, "users-service-api ", log.LstdFlags)

  err := godotenv.Load()

  if err != nil {
    l.Fatal("Error loading .env file")
  }

	mongoURI := flag.String("mongoURI", "mongodb://root:password@localhost:27017", "Database hostname url")
	serverAddr := flag.String("serverAddr", "localhost", "Network address")
	serverPort := flag.Int("serverPort", 9090, "Port")

	co := options.Client().ApplyURI(*mongoURI)



	// establish database connection

	client, err := mongo.NewClient(co)

	if err != nil {
		l.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		l.Fatal(err)
	}


	defer func ()  {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	fmt.Println(err)


	serverURI := fmt.Sprintf("%s:%d", *serverAddr, *serverPort)

	fmt.Println("Database connected")

	sm := mux.NewRouter()


	getRouter := sm.Methods(http.MethodGet).Subrouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	putRouter := sm.Methods(http.MethodPut).Subrouter()


	uh := handlers.Newusers(l, client.Database("users-service-db").Collection("users"))

	getRouter.HandleFunc("/", Welcome)
	getRouter.HandleFunc("/users", uh.ListUsers)



	signUp := postRouter.PathPrefix("/signup").Subrouter()

	signUp.HandleFunc("/me", uh.Signup)
	signUp.Use(uh.ValidateUserData)

	postRouter.HandleFunc("/login", uh.Login)



	putRouter.HandleFunc("/users/{id:[0-9]+}", uh.UpdateUsers)
	putRouter.Use(uh.ValidateUserData)

	s := http.Server{
		Addr: serverURI,
		Handler: sm,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	
	go func ()  {
		err = s.ListenAndServe()
		if err != nil {
			l.Fatal("Server failed to start")
		}
	} ()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <- sigChan

	l.Println("Received kill signal", sig)
	
	ct, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	s.Shutdown(ct)
}
