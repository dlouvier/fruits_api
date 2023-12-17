package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestFruitsApi(t *testing.T) {
	t.Run("should return all fruits in JSON format", func(t *testing.T) {
		fruitsApi := New(fiber.New(), map[string]Fruit{})

		data := map[string]Fruit{"idOne": {"idOne", "a", "a"}}
		fruitsApi.data = data

		r := httptest.NewRequest("GET", "/api/fruits", nil)
		resp, _ := fruitsApi.app.Test(r)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		expected := convertMapToSlice(data)

		actual := []Fruit{}

		if assert.NoError(t, json.Unmarshal(body, &actual)) {
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, expected, actual)
		}

	})

	t.Run("should return a specific fruit in JSON format", func(t *testing.T) {
		data := map[string]Fruit{
			"CvNlZ0F3": {"CvNlZ0F3", "apple", "green"},
			"NoxN7Qm6": {"NoxN7Qm6", "banana", "yellow"},
			"h29Z90Oa": {"h29Z90Oa", "orange", "orange"},
		}

		fruitsApi := New(fiber.New(), data)

		r := httptest.NewRequest("GET", "/api/fruits/h29Z90Oa", nil)
		resp, _ := fruitsApi.app.Test(r)

		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)

		expected := Fruit{"h29Z90Oa", "orange", "orange"}

		actual := Fruit{}

		if assert.NoError(t, json.Unmarshal(body, &actual)) {
			assert.Equal(t, 200, resp.StatusCode)
			assert.Equal(t, expected, actual)
		}

		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("should add a fruit (multiple times) while keeping persistency across requests", func(t *testing.T) {
		fruitsApi := New(fiber.New(), map[string]Fruit{})

		fruitsToAdd := []Fruit{
			{"CvNlZ0F3", "apple", "green"},
			{"", "banana", "yellow"}, // intentionally without Id
		}
		expected := []Fruit{}

		for it, fruit := range fruitsToAdd {
			jsonData, _ := json.Marshal(fruit)

			r := httptest.NewRequest("POST", "/api/fruits", bytes.NewBuffer(jsonData))
			r.Header.Set("Content-Type", "application/json")

			resp, _ := fruitsApi.app.Test(r)
			assert.Equal(t, 200, resp.StatusCode)

			// recover the id for later comparison
			fruitId, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			expected = append(expected, Fruit{string(fruitId), fruit.Fruit, fruit.Color})

			t.Run(fmt.Sprintf("current size should be %d", it+1), func(t *testing.T) {
				r = httptest.NewRequest("GET", "/api/fruits", nil)
				resp, _ = fruitsApi.app.Test(r)

				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				fruits := []Fruit{}
				if assert.NoError(t, json.Unmarshal(body, &fruits)) {
					assert.Equal(t, it+1, len(fruits))
				}
			})
		}

		t.Run("should return contain both fruits", func(t *testing.T) {
			r := httptest.NewRequest("GET", "/api/fruits", nil)
			resp, _ := fruitsApi.app.Test(r)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			fruits := []Fruit{}
			if assert.NoError(t, json.Unmarshal(body, &fruits)) {
				assert.Contains(t, fruits, expected[0])
				assert.Contains(t, fruits, expected[1])
			}
		})

	})

	t.Run("should return results from search", func(t *testing.T) {
		data := map[string]Fruit{
			"CvNlZ0F3": {"CvNlZ0F3", "apple", "red"},
			"h29Z90Oa": {"h29Z90Oa", "apple", "Green"},
			"yJwk2_z":  {"yJwk2_z", "watermelon", "green"},
			"Ki2-xd":   {"Ki2-xd", "banana", "yellow"},
		}
		fruitsApi := New(fiber.New(), data)

		t.Run("should return all fruits with green color", func(t *testing.T) {
			jsonData, _ := json.Marshal(Fruit{"", "", "green"})

			r := httptest.NewRequest("POST", "/api/fruits/search", bytes.NewBuffer(jsonData))
			r.Header.Set("Content-Type", "application/json")

			resp, _ := fruitsApi.app.Test(r)
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			expected := []Fruit{{"h29Z90Oa", "apple", "Green"}, {"yJwk2_z", "watermelon", "green"}}
			actual := []Fruit{}
			if assert.NoError(t, json.Unmarshal(body, &actual)) {
				assert.Contains(t, actual, expected[0])
				assert.Contains(t, actual, expected[1])
			}

		})

		t.Run("should return all fruits apple", func(t *testing.T) {
			jsonData, _ := json.Marshal(Fruit{"", "Apple", ""})

			r := httptest.NewRequest("POST", "/api/fruits/search", bytes.NewBuffer(jsonData))
			r.Header.Set("Content-Type", "application/json")

			resp, _ := fruitsApi.app.Test(r)
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			expected := []Fruit{{"CvNlZ0F3", "apple", "red"}, {"h29Z90Oa", "apple", "Green"}}
			actual := []Fruit{}
			if assert.NoError(t, json.Unmarshal(body, &actual)) {
				assert.Equal(t, 200, resp.StatusCode)
				assert.Contains(t, actual, expected[0])
				assert.Contains(t, actual, expected[1])
			}

		})

		t.Run("should return fruit with ID Ki2-xd and color red", func(t *testing.T) {
			jsonData, _ := json.Marshal(Fruit{"Ki2-xd", "", "red"})

			r := httptest.NewRequest("POST", "/api/fruits/search", bytes.NewBuffer(jsonData))
			r.Header.Set("Content-Type", "application/json")

			resp, _ := fruitsApi.app.Test(r)
			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			expected := []Fruit{{"Ki2-xd", "banana", "yellow"}, {"CvNlZ0F3", "apple", "red"}}
			actual := []Fruit{}
			if assert.NoError(t, json.Unmarshal(body, &actual)) {
				assert.Equal(t, 200, resp.StatusCode)
				assert.Contains(t, actual, expected[0])
				assert.Contains(t, actual, expected[1])
			}
		})

		t.Run("should return status 204 no content", func(t *testing.T) {
			jsonData, _ := json.Marshal(Fruit{"", "", "orange"})

			r := httptest.NewRequest("POST", "/api/fruits/search", bytes.NewBuffer(jsonData))
			r.Header.Set("Content-Type", "application/json")

			resp, _ := fruitsApi.app.Test(r)
			assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		})

	})
}
