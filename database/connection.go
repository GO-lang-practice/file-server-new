package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

var Client *mongo.Client
var Database *mongo.Database

func Connect() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)

	// Create a context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Test the connection with a ping
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	Client = client

	// Initialize the Database variable
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "FileServer"
	}
	Database = client.Database(dbName)

	log.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func GetCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("DATABASE_NAME")
	if dbName == "" {
		dbName = "FileServer"
	}
	return Client.Database(dbName).Collection(collectionName)
}
