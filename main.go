package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

var collection *mongo.Collection

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOpt := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOpt)

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging to database", err)
	}
	fmt.Println("Connected to database")
	fmt.Println(client)
	fmt.Println(collection)
	fmt.Println(clientOpt)
	// collection = client.Database("todo").Collection("todo")
}
