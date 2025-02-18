package main

import (
	"fmt"
	"net/http"
	"strings"

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
    newsService, err := services.LoadMicroservice("ANIME_NEWS_API")
    if err != nil {
		panic(err)
	}
    forumService, err := services.LoadMicroservice("ANIME_FORUM_API")
    if err != nil {
		panic(err)
	}
	// Create a Gin router instance
	router := gin.Default()

    //! Anime Section
	router.GET("/anime", func(c *gin.Context) {
		body, statusCode, err := services.FetchDataFromMicroservice(animeService, "anime")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
	})

    router.GET("/anime/:id", func(c *gin.Context) {
		id := c.Param("id") // Get the 'id' parameter from the route 
		body, statusCode, err := services.FetchDataFromMicroservice(animeService, "anime/"+id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
    }) 

        // Define routes for serving images and videos
    router.GET("/uploads/videos/*filename", func(c *gin.Context) {
        // Extract filename and build endpoint path
        filename := strings.TrimPrefix(c.Param("filename"), "/")
        endpoint := fmt.Sprintf("uploads/videos/%s", filename)
        
        // Forward request to microservice
        body, statusCode, contentType, err := services.GetFileFromMicroservice(animeService, endpoint)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        // Return response from microservice
        c.Data(statusCode, contentType, body)
    })

    router.GET("/uploads/images/*filename", func(c *gin.Context) {
        // Extract filename and build endpoint path
        filename := strings.TrimPrefix(c.Param("filename"), "/")
        endpoint := fmt.Sprintf("uploads/images/%s", filename)
        
        // Forward request to microservice
        body, statusCode, contentType, err := services.GetFileFromMicroservice(animeService, endpoint)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        // Return response from microservice
        c.Data(statusCode, contentType, body)
    })

    //! News Section
    router.GET("/news", func(c *gin.Context) {
		body, statusCode, err := services.FetchDataFromMicroservice(newsService, "news")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
	})
    router.GET("/news/:id", func(c *gin.Context) {
		id := c.Param("id") // Get the 'id' parameter from the route 
		body, statusCode, err := services.FetchDataFromMicroservice(newsService, "news/"+id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
    })

    //! Forum Section
    router.GET("/threads", func(c *gin.Context) {
		body, statusCode, err := services.FetchDataFromMicroservice(forumService, "threads")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
	})
    router.GET("/posts/:thread_id", func(c *gin.Context) {
		id := c.Param("thread_id") // Get the 'id' parameter from the route 
		body, statusCode, err := services.FetchDataFromMicroservice(forumService, "posts/"+id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
    })
    router.GET("/threads/anime/:anime_id", func(c *gin.Context) {
		id := c.Param("anime_id") // Get the 'id' parameter from the route 
		body, statusCode, err := services.FetchDataFromMicroservice(forumService, "threads/anime/"+id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Data(statusCode, "application/json", body)
    })

    router.POST("/register", func(c *gin.Context) {
        // Extract form data
        formData := map[string]string{
            "username": c.PostForm("username"),
            "password": c.PostForm("password"),
            "profile_url": c.PostForm("profile_url"),
        }

        // Ensure required fields exist
        if formData["username"] == "" || formData["password"] == "" || formData["profile_url"] == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
            return
        }

        // Send to microservice
        body, statusCode, err := services.PostDataToMicroservice(forumService, "/register", formData)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Respond with the microservice response
        c.Data(statusCode, "application/json", body)
    })

    router.POST("/login", func(c *gin.Context) {
        // Extract form data
        formData := map[string]string{
            "username": c.PostForm("username"),
            "password": c.PostForm("password"), 
        }

        // Ensure required fields exist
        if formData["username"] == "" || formData["password"] == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
            return
        }

        // Send to microservice
        body, statusCode, err := services.PostDataToMicroservice(forumService, "/login", formData)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Respond with the microservice response
        c.Data(statusCode, "application/json", body)
    })

    router.POST("/user", func(c *gin.Context) {
        // Extract form data
        formData := map[string]string{
            "session_token": c.PostForm("session_token"), 
        }

        // Ensure required fields exist
        if formData["session_token"] == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
            return
        }

        // Send to microservice
        body, statusCode, err := services.PostDataToMicroservice(forumService, "/user", formData)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Respond with the microservice response
        c.Data(statusCode, "application/json", body)
    })

    router.POST("/thread", func(c *gin.Context) {
        // Extract form data
        formData := map[string]string{
            "title": c.PostForm("title"),
            "author_id": c.PostForm("author_id"),
            "anime_id": c.PostForm("anime_id"),
        }

        // Ensure required fields exist
        if formData["title"] == "" || formData["author_id"] == "" || formData["anime_id"] == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
            return
        }

        // Send to microservice
        body, statusCode, err := services.PostDataToMicroservice(forumService, "/thread", formData)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Respond with the microservice response
        c.Data(statusCode, "application/json", body)
    })

    router.POST("/post", func(c *gin.Context) {
        // Extract form data
        formData := map[string]string{
            "content": c.PostForm("content"),
            "author_id": c.PostForm("author_id"),
            "thread_id": c.PostForm("thread_id"),
        }

        // Ensure required fields exist
        if formData["content"] == "" || formData["author_id"] == "" || formData["thread_id"] == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
            return
        }

        // Send to microservice
        body, statusCode, err := services.PostDataToMicroservice(forumService, "/post", formData)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Respond with the microservice response
        c.Data(statusCode, "application/json", body)
    })
	// Start the server on port 8080
	router.Run(":8080")
}

