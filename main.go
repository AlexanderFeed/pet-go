package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
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

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging to database", err)
	}
	fmt.Println("Connected to database")
	// The line `// collection = client.Database("todo").Collection("todo")` is a commented-out line of
	// code in Go. It is not being executed as it is preceded by `//`, which indicates a single-line
	// comment in Go.
	collection = client.Database("golang_db").Collection("todo")
	app := fiber.New()
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen(":" + port))
}

func getTodos(c *fiber.Ctx) error {
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		c.Status(500).SendString(err.Error())
		return err
	}
	defer cursor.Close(context.Background()) /////???
	var todos []Todo
	for cursor.Next(context.Background()) {
		var todo Todo
		cursor.Decode(&todo)
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}

func createTodos(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		c.Status(400).SendString(err.Error())
		return err
	}
	collection.InsertOne(context.Background(), todo)
	return c.JSON(todo)
}

func updateTodos(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(400).SendString(err.Error())
		return err
	}
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		c.Status(400).SendString(err.Error())
		return err
	}
	todo.Id = id
	collection.UpdateOne(context.Background(), bson.M{"id": id}, bson.M{"$set": todo})
	return c.JSON(todo)
}
func deleteTodos(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Status(400).SendString(err.Error())
		return err
	}
	collection.DeleteOne(context.Background(), bson.M{"id": id})
	return c.SendStatus(fiber.StatusOK)
}
