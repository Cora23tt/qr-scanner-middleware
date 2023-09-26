package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type Scan struct {
	Item    string `json:"item"`
	OrgCode string `json:"org_code"`
}

var (
	scans []Scan
	mutex sync.Mutex
)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/scans", func(c *gin.Context) {
		var newScans []string
		if err := c.BindJSON(&newScans); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		mutex.Lock()
		for _, newItem := range newScans {
			// Remove the last 3 symbols as the organization code
			orgCode := newItem[len(newItem)-3:]
			item := newItem[:len(newItem)-3]

			// Add the item and organization code to the scans array
			scans = append(scans, Scan{Item: item, OrgCode: orgCode})
		}
		mutex.Unlock()
		c.JSON(http.StatusCreated, gin.H{"message": "Scans added successfully"})
	})
	r.GET("/scans/:orgCode", func(c *gin.Context) {
		requestedOrgCode := c.Param("orgCode")
		mutex.Lock()
		var filteredScans []Scan
		for _, scan := range scans {
			if scan.OrgCode == requestedOrgCode {
				filteredScans = append(filteredScans, scan)
			}
		}
		mutex.Unlock()
		c.JSON(http.StatusOK, filteredScans)
	})
	r.DELETE("/scans/:orgCode", func(c *gin.Context) {
		orgCodeToDelete := c.Param("orgCode")
		mutex.Lock()
		// Create a new slice to store scans without the specified orgCode
		var newScans []Scan
		for _, scan := range scans {
			if scan.OrgCode != orgCodeToDelete {
				newScans = append(newScans, scan)
			}
		}
		// Update the scans array
		scans = newScans
		mutex.Unlock()
		c.JSON(http.StatusOK, gin.H{"message": "Scans cleared successfully for OrgCode: " + orgCodeToDelete})
	})
	r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH,OPTIONS,GET,PUT,DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
