package main

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	mongodbEnv = "MONGODB_URI"
	db         = "content_consolidation_db"
	collection = "content"
)

type Content struct {
	ID          bson.ObjectID `bson:"_id,omitempty"`
	Location    string
	Description string
	UpdatedAt   time.Time `bson:"updatedAt"`
}

type Repo struct {
	client *mongo.Client
}

func (r *Repo) Init() error {
	var err error
	clientOptions := options.Client().ApplyURI(os.Getenv(mongodbEnv))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r.client, err = mongo.Connect(clientOptions)
	if err != nil {
		return err
	}

	err = r.client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) SaveContentEntry(content Content) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new entry in the db
	content.UpdatedAt = time.Now()

	filter := bson.M{"location": content.Location}
	update := bson.M{"$set": content}

	opts := options.UpdateOne().SetUpsert(true)

	coll := r.client.Database(db).Collection(collection)
	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetContent() ([]Content, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set the collection
	coll := r.client.Database(db).Collection(collection)

	// Read from the collection
	cursor, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var content []Content
	if err = cursor.All(ctx, &content); err != nil {
		return nil, err
	}

	return content, nil
}

/*
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
*/
