package all

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler

	_ "github.com/dlouvier/fruits-api/src/docs"
)

var dataFilename = "./fruits-api-data.json"

type Fruit struct {
	Id    string `json:"id"`
	Fruit string `json:"fruit" example:"apple"`
	Color string `json:"color" example:"red"`
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
// @Success 200 {string} main.Fruit.Id
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

// @Summary Search a fruit
// @Description Search for a fruit based in a Payload
// @Accept json
// @Produce json
// @Param request body main.Fruit true "query params"
// @Success 200 {object} []main.Fruit
// @Failure 204
// @Router /api/fruits/search [post]
func (fa *FruitsApi) searchFruit(c *fiber.Ctx) error {
	searchFruit := Fruit{}
	if err := c.BodyParser(&searchFruit); err != nil {
		return err
	}

	results := []Fruit{}

	if len(searchFruit.Id) > 0 {
		if fruit, exists := fa.data[searchFruit.Id]; exists {
			results = append(results, fruit)
		}
	}

	if len(searchFruit.Color) > 0 {
		for _, f := range fa.data {
			if strings.EqualFold(f.Color, searchFruit.Color) {
				results = append(results, f)
			}
		}
	}

	if len(searchFruit.Fruit) > 0 {
		for _, f := range fa.data {
			if strings.EqualFold(f.Fruit, searchFruit.Fruit) {
				results = append(results, f)
			}
		}
	}

	if len(results) > 0 {
		return c.Status(fiber.StatusOK).JSON(results)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
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
	apiFruits.Post("/search", fa.searchFruit)

	// swagger
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	return fa
}

func Save(data map[string]Fruit) {
	if len(data) > 0 {
		var buf bytes.Buffer

		encoder := json.NewEncoder(&buf)
		encoder.SetIndent("", "  ")
		encoder.SetEscapeHTML(false)

		err := encoder.Encode(data)
		if err != nil {
			log.Panic("Error encoding object")
		}

		file, _ := os.Create(dataFilename)
		err = os.WriteFile(file.Name(), buf.Bytes(), 0600)
		if err != nil {
			log.Panic("Error writing file")
		}

		log.Printf("Data saved to file '%s'", dataFilename)
		log.Printf("Saved data contains %d fruits", len(data))
	} else {
		log.Println("There were not fruits in the API. Not creating file")
	}

}

func Load() map[string]Fruit {
	file, err := os.Open(dataFilename)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			log.Println("Not found file with data")
			return map[string]Fruit{}
		} else {
			log.Panicln("Error opening file:", err)
		}
	}
	defer file.Close()

	byte, _ := io.ReadAll(file)

	data := map[string]Fruit{}
	err = json.Unmarshal(byte, &data)
	if err != nil {
		log.Panicln("Error opening file:", err)
	}

	return data
}
