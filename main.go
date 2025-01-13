package main

import (
	"io"
	"net/http"

	"anime-d-verse/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Load microservices
	animeService, err := services.LoadMicroservice("ANIME_CONTENT_API")
	if err != nil {
		panic(err)
	}

	// Create a Gin router instance
	router := gin.Default()

	// Define GET route for anime
	router.GET("/anime", func(c *gin.Context) {
		body, statusCode, err := services.FetchDataFromMicroservice(animeService, "anime")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
	})

	// Define POST route for anime
	router.POST("/anime", func(c *gin.Context) {
		// Read JSON payload from the request body
		jsonData, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		body, statusCode, err := services.PostDataToMicroservice(animeService, "anime", jsonData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
	})

	// Start the server on port 8080
	router.Run(":8080")
}
