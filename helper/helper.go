package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongodb-docker/models"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb://mongodb:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB...")

	collection := client.Database("bookshop").Collection("books")

	return collection
}

func GetError(err error, w http.ResponseWriter, statusCode int) {
	log.Println(err.Error())

	response := models.ErrorResponse{
		ErrorMessage: err.Error(),
		StatusCode:   statusCode,
	}

	message, _ := json.Marshal(response)

	w.WriteHeader(response.StatusCode)
	w.Write(message)
}

