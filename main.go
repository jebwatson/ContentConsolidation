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

const mongodbEnv = "MONGODB_URI"

var mongoClient *mongo.Client

type Website struct {
	ID          bson.ObjectID `bson:"_id"`
	Site        string
	Description string
}

func main() {
	connectToMongo(os.Getenv(mongodbEnv))
	defer closeMongoDB()
	results := ReadContentEntry()
	UpdateContentEntry(results[0])
	results = ReadContentEntry()
	DeleteContentEntry(results[0])
}

func CreateContentEntry() {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	db := mongoClient.Database("content_consolidation_db")
	coll := db.Collection("websites")
	doc := Website{Site: "https://google.com", Description: "A useful search engine"}

	result, insertErr := coll.InsertOne(ctx, doc)
	if insertErr != nil {
		log.Fatal("Could not insert record to collection: ", insertErr)
	}

	fmt.Println("Inserted document with _id: ", result.InsertedID)
}

func ReadContentEntry() []Website {
	// Create context
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

	fmt.Println("Found the following records:")
	for _, result := range results {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}

	return results
}

func UpdateContentEntry(website Website) {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	website.Site = "https://reddit.com"
	website.Description = "A hive of scum and villainy"
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")
	filter := bson.M{"_id": website.ID}
	update := bson.M{"$set": bson.M{
		"site":        website.Site,
		"description": website.Description,
	}}

	result, updateErr := coll.UpdateOne(ctx, filter, update)
	if updateErr != nil {
		log.Fatal("Could not update record ", website.ID, ": ", updateErr)
	}

	fmt.Println("Updated ", result.ModifiedCount, " records.")
}

func DeleteContentEntry(website Website) {
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")
	filter := bson.M{"_id": website.ID}

	result, deleteError := coll.DeleteOne(ctx, filter)
	if deleteError != nil {
		log.Fatal("Could not update record ", website.ID, ": ", deleteError)
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
