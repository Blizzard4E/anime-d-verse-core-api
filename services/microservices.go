package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
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

func PostDataToMicroservice(service Microservice, endpoint string, formData map[string]string) ([]byte, int, error) {
	if service.BaseURL == "" {
		return nil, 0, errors.New("microservice base URL is not set")
	}

	// Convert formData map to url.Values
	form := url.Values{}
	for key, value := range formData {
		form.Set(key, value)
	}
	encodedFormData := form.Encode()

	// Construct the full URL
	apiURL := service.BaseURL + "/" + endpoint

	// Create a new POST request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBufferString(encodedFormData))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
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


// GetFileFromMicroservice performs a GET request to fetch files from the specified microservice
func GetFileFromMicroservice(service Microservice, endpoint string) ([]byte, int, string, error) {
    if service.BaseURL == "" {
        return nil, 0, "", errors.New("microservice base URL is not set")
    }

    // Construct the full URL using the existing pattern
    apiURL := service.BaseURL + "/" + endpoint

    // Make the GET request
    resp, err := http.Get(apiURL)
    if err != nil {
        return nil, 0, "", err  // Return 0 status code for http errors
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, 0, "", err
    }

    // Get content type from response headers
    contentType := resp.Header.Get("Content-Type")
    return body, resp.StatusCode, contentType, nil
}