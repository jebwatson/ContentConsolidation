package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Website struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Site        string
	Description string
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

const mongodbEnv = "MONGODB_URI"

var mongoClient *mongo.Client

func CreateContentEntry(site string, description string) {
	// Connect to MongoDB
	uri := os.Getenv(mongodbEnv)
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	connectToMongo(uri)
	defer closeMongoDB() // Create context

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")
	doc := Website{Site: site, Description: description, CreatedAt: time.Now()}

	result, insertErr := coll.InsertOne(ctx, doc)
	if insertErr != nil {
		log.Fatal("Could not insert record to collection: ", insertErr)
	}

	fmt.Println("Inserted document with _id: ", result.InsertedID)
}

func ReadContentEntry() []Website {
	// Connect to MongoDB
	uri := os.Getenv(mongodbEnv)
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	connectToMongo(uri)
	defer closeMongoDB() // Create context

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set the collection
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")

	// Read from the collection
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Could not read records from MongoDB: ", err)
	}

	var results []Website
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal("Could not parse results from MongoDB: ", err)
	}

	return results
}

func UpdateContentEntry(website Website) {
	// Connect to MongoDB
	uri := os.Getenv(mongodbEnv)
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	connectToMongo(uri)
	defer closeMongoDB()

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")
	filter := bson.M{"_id": website.ID}
	update := bson.M{"$set": bson.M{
		"site":        website.Site,
		"description": website.Description,
		"UpdatedAt":   time.Now(),
	}}

	result, updateErr := coll.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		log.Fatal("Could not update record ", website.ID, ": ", updateErr)
	}

	fmt.Println("Updated ", result.ModifiedCount, " records.")
}

func DeleteContentEntry(id bson.ObjectID) {
	// Connect to MongoDB
	uri := os.Getenv(mongodbEnv)
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not set")
	}
	connectToMongo(uri)
	defer closeMongoDB()

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")
	filter := bson.M{"_id": id}

	result, deleteError := coll.DeleteOne(ctx, filter)
	if deleteError != nil {
		log.Fatal("Could not delete record ", id, ": ", deleteError)
	}

	fmt.Println("Deleted ", result.DeletedCount, " records.")
}

func connectToMongo(uri string) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		fmt.Println(uri)
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
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
