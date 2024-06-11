package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBInstance()

func DBInstance() *mongo.Client {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
	}
	url := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		fmt.Println(err.Error())
	}
	return client
}

func OpenCollection(conn *mongo.Client, collection string) *mongo.Collection {
	coll := conn.Database("File-Sharing").Collection(collection)
	return coll
}
