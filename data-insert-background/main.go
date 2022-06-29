package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

var ctx = context.TODO()

func main() {

	connectionURI := os.Args[1]
	concurrentExecutions, _ := strconv.Atoi(os.Args[2])

	ch := make(chan string)

	for i := 0; i < concurrentExecutions; i++ {
		go insertData(connectionURI, ch)
	}

	for {
		go insertData(<-ch, ch)
	}

}

func insertData(connectionURI string, ch chan string) {
	// this is the primary host for the replicaset
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
	// so we don't flood the DB with connections
	time.Sleep(300 * time.Millisecond)

	log.Printf("inserting data query...")

	a := FakeData{}
	err = faker.FakeData(&a)
	if err != nil {
		fmt.Println(err)
	}

	result, err := mongo_collection.InsertOne(ctx, a)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("inserted ID: %v\n", result.InsertedID)
	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	ch <- connectionURI

}
