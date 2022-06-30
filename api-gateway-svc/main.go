package main

import (
	"api-gateway-svc/config"
	"api-gateway-svc/services/account"
	"api-gateway-svc/services/order"
	"api-gateway-svc/services/product"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)





func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatal("Failed to load config", err)
	}

	sm := mux.NewRouter()

	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Context-Type", "Application/json")
		fmt.Fprint(w, "Ok")
	})


	sm.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Here we go gateway"))
	})



	http.Handle("/", sm)


	account.RegisterService(sm, c)
	product.RegisterService(sm, c)
	order.RegisterService(sm, c)

	http.ListenAndServe(c.Port, sm)

	s := http.Server {
		Addr: c.Port,
		IdleTimeout: 120 * time.Second,
		ReadTimeout: 1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}



	go func ()  {
		err = 	s.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to server")
		}
	}()


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	<-sigChan
	
	fmt.Println("Received signal to kill server")

	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()
	
	s.Shutdown(ctx)
}