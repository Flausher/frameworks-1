package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
)

const baseURL = "http://127.0.0.1:8080"

func TestAPI_E2E(t *testing.T) {
	t.Run("GET /api/users - Success", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/users", baseURL))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got: %d", resp.StatusCode)
		}
	})

	t.Run("GET /api/users/{id} - Not Found", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/api/users/9999", baseURL))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected 404 Not Found, got: %d", resp.StatusCode)
		}
	})

	t.Run("POST /api/users - Invalid Data", func(t *testing.T) {
		badPayload := []byte(`{"email": "bad@mail.com", "name": "", "age": -5}`)
		resp, err := http.Post(fmt.Sprintf("%s/api/users", baseURL), "application/json", bytes.NewBuffer(badPayload))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected 400 Bad Request, got: %d", resp.StatusCode)
		}
	})

	var createdID int
	t.Run("POST /api/users - Success", func(t *testing.T) {
		payload := []byte(`{
			"email": "test@example.com",
			"password": "secure123",
			"name": "Test User",
			"age": 25
		}`)
		resp, err := http.Post(fmt.Sprintf("%s/api/users", baseURL), "application/json", bytes.NewBuffer(payload))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected 200 or 201, got: %d", resp.StatusCode)
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &result); err == nil {
			if idFloat, ok := result["id"].(float64); ok {
				createdID = int(idFloat)
			}
		}
	})

	t.Run("GET /api/users/{id} - Success", func(t *testing.T) {
		if createdID == 0 {
			t.Skip("Skipping because we didn't get an ID from the POST test")
		}

		resp, err := http.Get(fmt.Sprintf("%s/api/users/%d", baseURL, createdID))
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200 OK, got: %d", resp.StatusCode)
		}
	})
}
