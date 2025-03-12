package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func LoadJsonFromResultFile(filepath string) ([]gin.H, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("error during the openning or the file : %v", err)
	}
	defer file.Close()

	var data []gin.H
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("error during json's decodage : %v", err)
	}

	return data, nil
}

func getUploadsName() string {
	filepath := "uploads"
	files, err := os.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		return file.Name()
	}

	return ""
}

// Route /analyse
func StartAnalyse(c *gin.Context) {
	uploadsName := "uploads/" + getUploadsName() + "/"
	if uploadsName == "uploads//" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid upload found"})
		return
	}

	cmd := exec.Command("./exec.sh", uploadsName)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		log.Fatal("Start analyse request err:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Analysis failed"})
		return
	}

	newJson, err := LoadJsonFromResultFile("result.json")
	if err != nil {
		log.Fatal("Error: result.json doesn't exist", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "result.json not found"})
		return
	}

	c.JSON(http.StatusOK, newJson)
}
