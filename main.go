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
