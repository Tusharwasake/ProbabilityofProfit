package main

import (
	"fmt"
	"log"
	"os"
	"pop-calculator/controller"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables (optional - only for SERVER_PORT)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found - using defaults")
	}

	app := gin.Default()
	
	// Simple health check endpoint
	app.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
			"service": "Probability of Profit Calculator",
		})
	})

	// Main PoP calculation endpoint
	app.POST("/pop", controller.CalculatePoP)
	
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	app.Run(":" + port)
}
