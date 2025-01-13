package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
)

// Microservice represents a predefined microservice with its base URL
type Microservice struct {
	BaseURL string
}

// LoadMicroservice creates a Microservice instance with its base URL loaded from an environment variable
func LoadMicroservice(envVar string) (Microservice, error) {
	baseURL := os.Getenv(envVar)
	if baseURL == "" {
		return Microservice{}, errors.New("environment variable " + envVar + " is not set")
	}
	return Microservice{BaseURL: baseURL}, nil
}

// FetchDataFromMicroservice performs a GET request to the specified microservice and endpoint
func FetchDataFromMicroservice(service Microservice, endpoint string) ([]byte, int, error) {
	if service.BaseURL == "" {
		return nil, 0, errors.New("microservice base URL is not set")
	}

	// Construct the full URL
	apiURL := service.BaseURL + "/" + endpoint

	// Make the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}

// PostDataToMicroservice performs a POST request to the specified microservice and endpoint
func PostDataToMicroservice(service Microservice, endpoint string, jsonData []byte) ([]byte, int, error) {
	if service.BaseURL == "" {
		return nil, 0, errors.New("microservice base URL is not set")
	}

	// Construct the full URL
	apiURL := service.BaseURL + "/" + endpoint

	// Create a new POST request with the JSON data
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, resp.StatusCode, nil
}
