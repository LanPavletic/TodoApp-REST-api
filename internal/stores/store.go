package stores

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MONGO_URI     = os.Getenv("MONGO_URI")
	DATABASE_NAME = os.Getenv("MONGO_DATABASE_NAME")
	JWT_SECRET    = os.Getenv("JWT_SECRET")
)

func InitMongoDB() *mongo.Client {
	opts := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func Disconnect(c *mongo.Client) {
	err := c.Disconnect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
