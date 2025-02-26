package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)
const temp_db = "./temp_db.txt"

var mu sync.Mutex

type shortenRequest struct {
	LongUrl string `json:"longUrl" binding:"required"` 
}

type shortenData struct {
	LongUrl string `json:"longUrl"`
}

func main() {
	router := gin.Default()

	router.GET("/", handleHome)
	router.GET("/urls", handleGetUrls)
	router.POST("/shorten", handleShorten)

	fmt.Println("Server running on http://localhost:8080")
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleHome(c *gin.Context) {
	c.String(http.StatusOK,"This is my homepage")
}

func handleGetUrls(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.Open(temp_db)
	if err != nil {
		c.IndentedJSON(500, gin.H{"error":"Failed to read stored URLs"})
		return
	}
	defer file.Close()
	
	//Read and Parse Files
	var urls []shortenData
	decoder := json.NewDecoder(file)
	for decoder.More() {
		var url shortenData
		if err := decoder.Decode(&url); err != nil {
			c.IndentedJSON(500, gin.H{"error": "Error reading URLs"})
			return
		}
		urls = append(urls,url)
	}

	c.JSON(http.StatusOK, gin.H{"URLs": urls})
}

func handleShorten(c *gin.Context) {
	var req shortenRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := saveUrl(req.LongUrl); err != nil {
		c.IndentedJSON(500, gin.H{"error":"Failed to save URL"})
		return
	}	

	log.Printf("Stored URL: %s\n", req.LongUrl)

	c.IndentedJSON(http.StatusCreated,  gin.H{
		"message": "received and appended the URL successfully",
		"longUrl": req.LongUrl,
	})
}

func saveUrl(longUrl string) error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(temp_db, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data := shortenData{LongUrl: longUrl}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_,err = file.WriteString(string(jsonData) + "\n")
	return err
}