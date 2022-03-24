package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/bxcodec/faker/v3"
)

type FakeData struct {
	Name            string `faker:"name"`
	UserName        string `faker:"username"`
	Email           string `faker:"email"`
	ReviewSubmitted string `faker:"timestamp"`
	Review          string `faker:"paragraph"`
	BusinessId      int    `faker:"boundary_start=1, boundary_end=10000"`
}

func main() {

	numberOfRecords, err := strconv.Atoi(os.Args[1])
	// this is for making it easier to test locally vs running in a container
	outPutDirectory := os.Args[2]

	jsonData := []FakeData{}
	for i := 0; i < numberOfRecords+1; i++ {
		a := FakeData{}
		err := faker.FakeData(&a)
		if err != nil {
			fmt.Println(err)
		}

		if i%1000 == 0 {
			fmt.Printf("generated record %d out of %d\n", i, numberOfRecords)
		}
		jsonData = append(jsonData, a)
		if i > 0 {
			if (i % 100000) == 0 {
				comment := fmt.Sprintf("writing out to file: record-%d.json", i)
				fmt.Println(comment)
				jsonBlob, _ := json.MarshalIndent(jsonData, "", " ")
				_ = ioutil.WriteFile(fmt.Sprintf("%s/record-%d.json", outPutDirectory, i), jsonBlob, 0644)
				jsonData = []FakeData{}
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Finished writing to file.")
}
