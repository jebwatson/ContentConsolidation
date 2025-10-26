package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongodb_env = "MONGODB_URI"

var mongoClient *mongo.Client

type Website struct {
	Uri         string
	Description string
}

func main() {
	connectToMongo(os.Getenv(mongodb_env))
	defer closeMongoDB()
	ReadContentEntry()
}

func CreateContentEntry() {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	db := mongoClient.Database("content_consolidation_db")
	coll := db.Collection("websites")
	doc := Website{Uri: "https://google.com", Description: "A useful search engine"}

	result, insertErr := coll.InsertOne(ctx, doc)
	if insertErr != nil {
		log.Fatal("Could not insert record to collection: %v", insertErr)
	}

	fmt.Println("Inserted document with _id: %v\n", result.InsertedID)
}

func ReadContentEntry() {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set the collection
	db := mongoClient.Database("content_consolidation_db")
	coll := db.Collection("websites")

	// Read from the collection
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Could not read records from MongoDB: %v", err)
	}

	var results []Website
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal("Could not parse results from MongoDB: %v", err)
	}

	fmt.Println("Found the following records:\n")
	for _, result := range results {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}
}

func connectToMongo(uri string) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(uri)
		log.Fatal("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	mongoClient = client
}

func closeMongoDB() {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			fmt.Println("Disconnected from MongoDB.")
		}
	}
}
