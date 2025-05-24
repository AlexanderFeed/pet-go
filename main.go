package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	Id        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello Worlds")

	app := fiber.New()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	todos := []Todo{}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"message": "Hello World",
			"todos":   todos,
		})
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} //{id: 0, completed: false, body: ""}

		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{
				"message": "Body cannot be empty",
			})
		}

		todo.Id = len(todos) + 1
		todos = append(todos, *todo)
		return c.Status(200).JSON(fiber.Map{
			"message": "Todo created",
			"todo":    todo,
		})
	})

	//update todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}
		for i, todo := range todos {
			if todo.Id == id {
				todos[i].Completed = !todo.Completed
				return c.Status(200).JSON(fiber.Map{
					"message": "Todo updated",
					"todo":    todo,
				})
			}
		}
		return c.Status(404).JSON(fiber.Map{
			"message": "Todo not found",
		})
	})

	//delete

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return err
		}
		for i, todo := range todos {
			if todo.Id == id {
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{
					"message": "Todo deleted",
					"todo":    todo,
				})
			}
		}
		return c.Status(404).JSON(fiber.Map{
			"message": "Todo not found",
		})
	})

	log.Fatal(app.Listen(":" + PORT))

}
