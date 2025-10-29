package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const mongodbEnv = "MONGODB_URI"

type Website struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Site        string
	Description string
	CreatedAt   time.Time `bson:"createdAt"`
	UpdatedAt   time.Time `bson:"updatedAt"`
}

type MongoError struct {
	Message string
}

var mongoClient *mongo.Client = nil

func (e *MongoError) Error() string {
	return string(e.Message)
}

func CreateContentEntry(site string, description string) (string, *MongoError) {
	err := connectToMongo()
	if err != nil {
		return "", err
	}
	defer disconnectFromMongo()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")
	doc := Website{Site: site, Description: description, CreatedAt: time.Now()}

	result, insertErr := coll.InsertOne(ctx, doc)
	if insertErr != nil {
		return "", &MongoError{"Could not create record: " + insertErr.Error()}
	}

	return fmt.Sprintf("Created document with _id: %s", result.InsertedID), nil
}

func ReadContentEntries() ([]Website, *MongoError) {
	connectToMongo()
	defer disconnectFromMongo()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set the collection
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")

	// Read from the collection
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, &MongoError{fmt.Sprintf("Could not read records from MongoDB: %s", err)}
	}

	var results []Website
	if err = cursor.All(ctx, &results); err != nil {
		return nil, &MongoError{fmt.Sprintf("Could not parse results from MongoDB: %s", err)}
	}

	return results, nil
}

func UpdateContentEntry(website Website) (string, *MongoError) {
	connectToMongo()
	defer disconnectFromMongo()

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
		return "", &MongoError{fmt.Sprintf("Could not update record %s: %s", website.ID, updateErr)}
	}

	return fmt.Sprintf("Updated %s records.", result.ModifiedCount), nil
}

func DeleteContentEntry(id bson.ObjectID) (string, *MongoError) {
	connectToMongo()
	defer disconnectFromMongo()

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	coll := mongoClient.Database("content_consolidation_db").Collection("websites")
	filter := bson.M{"_id": id}

	result, deleteError := coll.DeleteOne(ctx, filter)
	if deleteError != nil {
		return "", &MongoError{fmt.Sprintf("Could not delete record %s: %s", id, deleteError.Error())}
	}

	return fmt.Sprint("Deleted %v records.", result.DeletedCount), nil
}

func connectToMongo() (error *MongoError) {
	clientOptions := options.Client().ApplyURI(os.Getenv(mongodbEnv))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return &MongoError{"Failed to connect to MongoDB: " + err.Error()}
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return &MongoError{"Failed to ping MongoDB: " + err.Error()}
	}

	mongoClient = client

	return nil
}

func disconnectFromMongo() (error *MongoError) {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			return &MongoError{"Error disconnecting from MongoDB: %v" + err.Error()}
		}
	}

	return nil
}
