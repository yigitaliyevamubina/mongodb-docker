package main

import (
	"fmt"
	"log"
	"mongodb-docker/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/books", handlers.ListBooks).Methods("GET")          //get list
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")       //get single
	router.HandleFunc("/books", handlers.CreateBook).Methods("POST")        //create
	router.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")    //update
	router.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE") //delete

	fmt.Println("Listening on book-service:8080...")
	log.Fatal(http.ListenAndServe("book-service:8080", router))
}