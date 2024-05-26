// Package main is the entry point of the program.
// the function main() is going to be executed when the program runs.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client // create a client variable to store the connection to the database

func pingDatabase() error {
	// ping the database to confirm a successful connection
	var result bson.M
	err := client.Database("robots").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result)

	if err != nil {
		return fmt.Errorf("ping error: %v", err)
	}

	return nil
}

func connectToMongo(uri string) error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	var err error
	// initialize the client and connect to the database
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return fmt.Errorf("connection error: %v", err) // error strings must be in lowercase, %v is used to print the value of the error, the second argument is the error itself
	}

	// by including the if statement, we can check if the connection was successful because the funciton will return an error if it wasn't
	if err := pingDatabase(); err != nil {
		return fmt.Errorf("ping error: %v", err)
	}

	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return nil
}

func sayHi(w http.ResponseWriter, r *http.Request) {
	log.Println("Hello, World!")
	fmt.Fprintf(w, "Hello, World!")
}

// HANDLER FUNCTION
func getAllRobots() ([]bson.M, error) {
	collection := client.Database("robots").Collection("robots")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel() // when the function ends, cancel the context

	cursor, err := collection.Find(ctx, bson.M{}) // the bson.M{} is an empty filter, so it will return all documents in the collection, M stands for map
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer cursor.Close(ctx) // when the function ends, close the cursor

	var robots []bson.M
	if err = cursor.All(ctx, &robots); err != nil { // All is a helper function that decodes all documents in the cursor and appends them to the provided slice
		// &robots is a pointer to the slice of robots
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return robots, nil // return the slice of robots and a nil error
}

// HTTP HANDLER
func fetchAllRobots(w http.ResponseWriter, r *http.Request) {
	// Log request details
	log.Printf("Received \033[32m%s\033[0m request for \033[34m%s\033[0m from: %s", r.Method, r.URL.Path, r.RemoteAddr)

	robots, err := getAllRobots()

	if err != nil {
		log.Printf("error getting all the robots: %v", err)
		http.Error(w, "Failed to fetch robots data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w) is going to encode the robots slice and write it to the response writer, since we want to send the data back to the client, we use json
	if err := json.NewEncoder(w).Encode(robots); err != nil {
		log.Printf("Encode error: %v", err)
		http.Error(w, "Failed to encode robots data", http.StatusInternalServerError)
	}
}

func main() {
	if err := godotenv.Load(); err != nil { // load the .env file
		log.Fatal("Error loading .env file")
	}

	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI is not set in the .env file")
	}

	if err := connectToMongo(uri); err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	http.HandleFunc("/", sayHi)
	http.HandleFunc("/all-robots", fetchAllRobots)
	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "This is an error", http.StatusNotFound)
	})

	log.Println("Server is running on port 8008")
	if err := http.ListenAndServe(":8008", nil); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
