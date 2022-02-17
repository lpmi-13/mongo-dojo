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
	ReviewSubmitted string `faker:"timestamp"`
	Review          string `faker:"paragraph"`
}

var collection *mongo.Collection
var ctx = context.TODO()

func main() {

	numberOfRecords, err := strconv.Atoi(os.Args[1])
        connectionURI := os.Args[2]

	// this is the primary host for the replicaset running inside the containers
	clientOptions := options.Client().ApplyURI(connectionURI)
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
			fmt.Printf("inserted ID # %d: %v\n", i, result.InsertedID)
		}
	}

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}
