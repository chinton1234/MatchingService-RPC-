package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	// "io/ioutil"
	"log"
	"net/http"
)

// "context"
// "context"
// "fmt"
// "time"

// // "time"

// "matchingService/models"

// "go.mongodb.org/mongo-driver/bson"
// "go.mongodb.org/mongo-driver/bson/primitive"

type response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func CreateMatching(id string, participant string) (data string, err error) {

	postBody, _ := json.Marshal(map[string]string{
		"userId": participant,
	})
	responseBody := bytes.NewBuffer(postBody)

	url := "http://172.31.86.56:8082" + "/matching/" + id

	//Leverage Go's HTTP Post function to make request
	resp, err := http.Post(url, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body

	if err != nil {
		return
	}

	mat := response{}
	b, err := io.ReadAll(resp.Body)
	fmt.Println(string(b))
	json.Unmarshal(b, &mat)

	data = mat.Data["data"].(string)
	return
}

func DeleteMatching(id string) (data string, err error) {

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", "http://172.31.86.56:8082/matching/"+id, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if err != nil {
		return
	}

	b, err := io.ReadAll(resp.Body)
	fmt.Println(string(b))

	if resp.Status == "404 Not Found" {
		data = "Matching Id not found."
	} else {
		data = "deleted."
	}

	return
}
