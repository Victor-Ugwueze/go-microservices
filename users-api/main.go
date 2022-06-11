package users

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/victor-ugwueze/go-micro/users-api/handlers"
	"github.com/victor-ugwueze/go-micro/users-api/models"
	"github.com/gorilla/mux"
)



func Welcome(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(wr, "Welcome")
}

func main() {

	users := []*models.User{}
	l := log.New(os.Stdout, "prod-api ", log.LstdFlags)

	uh := handlers.Newusers(l, &models.Models{ Users: users })

	sm := mux.NewRouter()


	getRouter := sm.Methods(http.MethodGet).Subrouter()
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()

	getRouter.HandleFunc("/", Welcome)
	getRouter.HandleFunc("/users", uh.ListUsers)


	postRouter.HandleFunc("/users", uh.CreateUsers)
	postRouter.Use(uh.ValidateUserData)

	putRouter.HandleFunc("/users/{id:[0-9]+}", uh.UpdateUsers)
	putRouter.Use(uh.ValidateUserData)
	
	deleteRouter.HandleFunc("/users/{id:[0-9]+}", uh.DeleteUsers)

	http.ListenAndServe(":9090", sm)
}
