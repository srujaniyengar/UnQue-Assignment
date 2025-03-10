package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// feat: connect to MongoDB using local URI and set global DB variable
func ConnectDB() *mongo.Database {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// feat: ping the database to verify connection
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Error pinging:", err)
	}

	DB = client.Database("unque_db")
	return DB
}
