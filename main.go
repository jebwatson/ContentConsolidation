package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const mongodb_env = "MONGODB_URI"

type Website struct {
	Uri         string
	Description string
}

func main() {
	CreateContentEntry()
}

func CreateContentEntry() {
	uri := os.Getenv(mongodb_env)
	if uri == "" {
		log.Fatal("Set your 'MONGODB_URI' environment variable.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, connectErr := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if connectErr != nil {
		log.Fatal("Error connecting to MongoDB: %v", connectErr)
	}

	pingErr := client.Ping(ctx, readpref.Primary())
	if pingErr != nil {
		log.Fatal("Could not ping MongoDB: %v", pingErr)
	}

	fmt.Println("Successfully connected to MongoDB!")

	// Create a new entry in the db
	db := client.Database("content_consolidation_db")
	coll := db.Collection("websites")
	doc := Website{Uri: "https://google.com", Description: "A useful search engine"}

	result, insertErr := coll.InsertOne(ctx, doc)
	if insertErr != nil {
		log.Fatal("Could not insert record to collection: %v", insertErr)
	}

	fmt.Println("Inserted document with _id: %v\n", result.InsertedID)

	// Disconnect from MongoDB when the application exits
	defer func() {
		if disconnectErr := client.Disconnect(ctx); disconnectErr != nil {
			log.Fatal("Error disconnecting from MongoDB: %v", disconnectErr)
		}
		fmt.Println("Disconnected from MongoDB")
	}()
}
