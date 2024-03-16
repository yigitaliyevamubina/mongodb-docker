package handlers

import (
	"context"
	"log"
	"mongodb-docker/helper"
	"mongodb-docker/models"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection = helper.ConnectMongoDB()

// Insert book
func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		helper.GetError(err, w, http.StatusInternalServerError)
		return
	}

	result, err := collection.InsertOne(context.TODO(), book)
	if err != nil {
		helper.GetError(err, w, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// Get book by object id
func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var book models.Book

	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		helper.GetError(err, w, http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}
	err = collection.FindOne(context.TODO(), filter).Decode(&book)
	if err != nil {
		helper.GetError(err, w, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// Update book by object id
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		helper.GetError(err, w, http.StatusBadRequest)
		return
	}

	var book models.Book

	filter := bson.M{"_id": id}

	err = json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		helper.GetError(err, w, http.StatusBadRequest)
		return
	}

	updateReq := bson.D{
		{"$set", bson.D{
			{"title", book.Title},
			{"price", book.Price},
			{"author", bson.D{
				{"firstname", book.Author.FirstName},
				{"lastname", book.Author.LastName},
				{"age", book.Author.Age},
			}},
		}},
	}

	err = collection.FindOneAndUpdate(context.TODO(), filter, updateReq).Decode(&book)
	if err != nil {
		helper.GetError(err, w, http.StatusInternalServerError)
		return
	}

	book.ID = id

	json.NewEncoder(w).Encode(book)
}

// Delete book by object id
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var params = mux.Vars(r)

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		helper.GetError(err, w, http.StatusBadRequest)
		return
	}

	filter := bson.M{"_id": id}

	deleteResponse, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		helper.GetError(err, w, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deleteResponse)
}

// List books
func ListBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var books []models.Book

	sortBy := bson.D{{"price", 1}}
	sortOptions := options.Find().SetSort(sortBy)
	cursor, err := collection.Find(context.TODO(), bson.D{}, sortOptions)
	if err != nil {
		helper.GetError(err, w, http.StatusInternalServerError)
		return
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var book models.Book

		err = cursor.Decode(&book)
		if err != nil {
			log.Println(err)
			return
		}

		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err)
		return
	}

	json.NewEncoder(w).Encode(books)
}
