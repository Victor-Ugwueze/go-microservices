package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/streadway/amqp"
)



func Welcome(wr http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(wr, "Welcome")
}

type OrderStatus uint16

const (
	StatusPending OrderStatus = iota
	StatusFailed
	StatusSuccess
	StatusComplete
)


type Order struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Price float64					`bson:"price"`
	UserId int						`bson:"userId"`
	Status OrderStatus 				`bson:"status"`
}

func main() {


	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://myuser:password@localhost:27017/shopping"))

	if err != nil {
		log.Fatal(err)
	}

	conn, rbErr := amqp.Dial("amqp://myuser:mypassword@localhost:5672/")

	if rbErr != nil {
		log.Fatalf("%s: %s", "Failed to connect to RabbitMQ", err)
		panic(rbErr)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("%s: %s", "Failed to open a channel", err)
		panic(rbErr)
	}

	defer ch.Close()


	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
			log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	database := client.Database("shopping")
	ordersCollection := database.Collection("orders")

	order := Order {
		Price: 100.0,
		UserId: 1,
	}

	newOrder, err := ordersCollection.InsertOne(ctx, order)

	if err != nil {
		panic(err)
	}

	fmt.Println(newOrder)

	cursor, err := ordersCollection.Find(ctx, bson.M{})

	if err != nil {
		panic(err)
	}


	var orders []bson.M

	if err = cursor.All(ctx, &orders); err != nil {
    log.Fatal(err)
  }

	fmt.Println(orders)

	sm := mux.NewRouter()


	



	http.ListenAndServe(":9090", sm)
}
