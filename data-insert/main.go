package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bxcodec/faker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FakeData struct {
	Name            string `faker:"name"`
	UserName        string `faker:"username"`
	Email           string `faker:"email"`
	ReviewRating    int    `faker:"oneof: 1, 2, 3, 4, 5"`
	ReviewSubmitted string `faker:"timestamp"`
	Review          string `faker:"paragraph"`
}

var collection *mongo.Collection
var ctx = context.TODO()

func main() {

	numberOfRecords, err := strconv.Atoi(os.Args[1])

	// this is the primary host for the replicaset running inside the VM
	clientOptions := options.Client().ApplyURI("mongodb://192.168.42.102:27017/?replicaSet=dojo")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongo_collection := client.Database("userData").Collection("reviews")

	for i := 0; i < numberOfRecords; i++ {
		a := FakeData{}
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println(err)
		}

		result, err := mongo_collection.InsertOne(ctx, a)
		if err != nil {
			log.Fatal(err)
		}
		if i%1000 == 0 {
			fmt.Printf("inserted ID: %v\n", result.InsertedID)
		}
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
