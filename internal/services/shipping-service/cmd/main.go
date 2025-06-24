package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get port from environment or use default
	port := os.Getenv("SHIPPING_SERVICE_PORT")
	if port == "" {
		port = "8005"
	}

	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "shipping-service",
			"status":  "healthy",
		})
	})

	// Shipping endpoints
	r.POST("/shipping/assign", assignShipper)
	r.GET("/shipping/track/:tracking_id", trackShipment)
	r.PUT("/shipping/:id/status", updateShippingStatus)
	r.POST("/shipping/:id/rate", rateShipper)

	// Start server
	log.Printf("Shipping Service starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func assignShipper(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "Shipper assigned successfully",
	})
}

func trackShipment(c *gin.Context) {
	trackingID := c.Param("tracking_id")
	c.JSON(http.StatusOK, gin.H{
		"tracking_id": trackingID,
		"status":      "in_transit",
	})
}

func updateShippingStatus(c *gin.Context) {
	shipmentID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"shipment_id": shipmentID,
		"message":     "Shipping status updated successfully",
	})
}

func rateShipper(c *gin.Context) {
	shipmentID := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"shipment_id": shipmentID,
		"message":     "Shipper rated successfully",
	})
}
