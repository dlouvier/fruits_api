package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler

	_ "github.com/dlouvier/fruits-api/src/docs"
)

type Fruit struct {
	Id    string `json:"id" example:"rKdzQ4"`
	Fruit string `json:"fruit" example:"Apple"`
	Color string `json:"color" example:"Red"`
}

type FruitsApi struct {
	app  *fiber.App
	data map[string]Fruit
}

func convertMapToSlice(m map[string]Fruit) []Fruit {
	fruits := []Fruit{}
	for id, fruit := range m {
		fruits = append(fruits, Fruit{id, fruit.Fruit, fruit.Color})
	}
	return fruits
}

func generateId() (string, error) {
	random := make([]byte, 6)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(random), nil
}

// @Summary Get the list of all fruits
// @Description Get the list of all fruits
// @Produce json
// @Success 200 {object} []main.Fruit
// @Router /api/fruits/ [get]
func (fa *FruitsApi) returnAll(c *fiber.Ctx) error {
	return c.JSON(convertMapToSlice(fa.data))
}

// @Summary Get a fruit by its Id
// @Description Get a fruit by its Id
// @Produce json
// @Param id path string true "Fruit ID" id
// @Success 200 {object} main.Fruit
// @Failure 404
// @Router /api/fruits/{id} [get]
func (fa *FruitsApi) returnOne(c *fiber.Ctx) error {
	id := c.Params("id")
	if fruit, exists := fa.data[id]; exists {
		return c.JSON(Fruit{fruit.Id, fruit.Fruit, fruit.Color})
	} else {
		return c.Status(fiber.StatusNotFound).SendString(fmt.Sprintf("Fruit with ID: '%s' not found", id))
	}
}

// @Summary Add a fruit
// @Description If the Payload doesn't contain an ID, it will generate it and also it doesn't validate if fields are empty or not. If the ID already exists, it will return an error.
// @Accept json
// @Produce json
// @Param request body main.Fruit true "query params"
// @Success 200 {string} main.FruitsJson.Id
// @Router /api/fruits [post]
func (fa *FruitsApi) addFruit(c *fiber.Ctx) error {
	newFruit := Fruit{}
	if err := c.BodyParser(&newFruit); err != nil {
		return err
	}

	if len(newFruit.Id) == 0 {
		id, err := generateId()
		if err != nil {
			return err
		}
		newFruit.Id = id
	}

	if _, exists := fa.data[newFruit.Id]; !exists {
		fa.data[newFruit.Id] = newFruit
		return c.Status(fiber.StatusOK).SendString(newFruit.Id)
	}

	// There is a chance the ID will collision.
	// To avoid using bigger hashes (longer ID) nor complicate the
	// code with retries mechanism, it just simply throw an error.
	return c.Status(fiber.StatusInternalServerError).SendString("ID is duplicated, try again")
}

func New(app *fiber.App, data map[string]Fruit) *FruitsApi {
	fa := &FruitsApi{
		app:  app,
		data: data,
	}

	api := app.Group("/api")
	apiFruits := api.Group("/fruits")
	apiFruits.Get("/", fa.returnAll)
	apiFruits.Get("/:id", fa.returnOne)
	apiFruits.Post("/", fa.addFruit)

	// swagger
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	return fa
}

func main() {
	data := map[string]Fruit{}
	fruitsApi := New(fiber.New(), data)
	fruitsApi.app.Listen(":3000")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var serverShutdown sync.WaitGroup

	go func() {
		fmt.Println("Gracefully shutting down...")
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = fruitsApi.app.ShutdownWithTimeout(60 * time.Second)
	}()

	if err := fruitsApi.app.Listen(":3000"); err != nil {
		log.Panic(err)
	}

	serverShutdown.Wait()
}
