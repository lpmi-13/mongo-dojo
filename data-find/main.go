package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// type FakeData struct {
// Name            string `faker:"name"`
// UserName        string `faker:"username"`
// Email           string `faker:"email"`
// ReviewSubmitted string `faker:"timestamp"`
// Review          string `faker:"paragraph"`
// BusinessId      int    `faker:"boundary_start=1, boundary_end=10000"`
// }

var collection *mongo.Collection
var ctx = context.TODO()

func main() {

	connectionURI := os.Args[1]
	concurrentExecutions, _ := strconv.Atoi(os.Args[2])

	ch := make(chan string)

	for i := 0; i < concurrentExecutions; i++ {
		go sendRequest(connectionURI, ch)
	}

	for {
		go sendRequest(<-ch, ch)
	}

}

func sendRequest(connectionURI string, ch chan string) {
	// this is the primary host for the replicaset running inside the VM
	clientOptions := options.Client().ApplyURI(connectionURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	filter := bson.M{"username": "wpgPoKz"}
	mongo_collection := client.Database("userData").Collection("reviews")
	time.Sleep(300 * time.Millisecond)
	result, err := mongo_collection.Find(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we found:", result)

	err = client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connection to MongoDB closed.")

	ch <- connectionURI

}
