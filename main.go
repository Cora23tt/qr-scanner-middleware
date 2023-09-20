package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	scans []string
	mutex sync.Mutex
)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/scanns", func(c *gin.Context) {
		var newScans []string
		if err := c.BindJSON(&newScans); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mutex.Lock()
		scans = append(scans, newScans...)
		mutex.Unlock()
		c.JSON(http.StatusCreated, gin.H{"message": "Scan added successfully"})
	})
	r.GET("/scanns", func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		data := scans
		scans = nil
		c.JSON(http.StatusOK, data)
	})
	r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS,GET,PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
