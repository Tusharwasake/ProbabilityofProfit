package main

import (
	"fmt"
	"log"

	"pop-calculator/controller"
	"pop-calculator/firstock"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}


	// Initialize Firstock client
	if err := firstock.InitializeFromEnv(); err != nil {
		log.Printf("Firstock authentication failed: %v", err)
		log.Println("Using fallback calculations")
	}


	app := gin.Default()
	
	app.GET("/status", func(c *gin.Context) {
		status := "authenticated"
		if firstock.JKey == "" {
			status = "fallback"
		}
		
		c.JSON(200, gin.H{
			"status": "running",
			"auth":   status,
		})
	})
	
	app.POST("/pop", controller.CalculatePoP)
	
	fmt.Println("Server starting on http://localhost:8080")
	app.Run(":8080")
}
