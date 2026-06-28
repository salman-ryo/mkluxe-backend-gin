package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectMongo establishes a connection to MongoDB and returns:
// 1. The MongoDB client (connection)
// 2. The selected database instance
// 3. An error if anything goes wrong
func ConnectMongo(uri, dbName string) (*mongo.Client, *mongo.Database, error) {

	// Create a context with a 10-second timeout.
	// The context is passed to MongoDB operations so they don't
	// wait forever if the server is unreachable.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// defer means "run this right before the function exits."
	//
	// cancel() DOES NOT disconnect MongoDB.
	// It simply cleans up the timeout context by stopping its timer
	// and releasing any resources associated with it.
	//
	// We defer it so that no matter how this function exits
	// (success or error), the context is always cleaned up.
	defer cancel()

	// Create a configuration object for the MongoDB client.
	// ApplyURI() tells the driver which MongoDB server to connect to.
	clientOptions := options.Client().ApplyURI(uri)

	// Attempt to establish a connection using the timeout context.
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		// If the connection couldn't be created, return the error.
		return nil, nil, err
	}

	// Verify that the MongoDB server is actually reachable.
	//
	// mongo.Connect() creates the client, but Ping() confirms
	// that the server is alive and responding.
	if err := client.Ping(ctx, nil); err != nil {
		return nil, nil, err
	}

	// Log a success message so we know the connection worked.
	log.Printf("Connected to database: %s", dbName)

	// Return:
	// - the MongoDB client (used for future database operations)
	// - the specific database instance
	// - nil because there was no error
	return client, client.Database(dbName), nil
}
