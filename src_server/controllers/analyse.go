package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	analyser "test_capstone/src_analyser"
	parser "test_capstone/src_parser"

	"github.com/gin-gonic/gin"
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

func getUploadsName(userName string) string {
	filepath := "uploads/" + userName

	files, err := os.ReadDir(filepath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		return file.Name()
	}

	return ""
}

func StartAnalyse(userName string, c *gin.Context) {
	uploadFolder := getUploadsName(userName)
	uploadsPath := fmt.Sprintf("uploads/%s/%s/", userName, uploadFolder)

	if uploadFolder == "" || userName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid upload found"})
		return
	}

	listOfFiles := parser.ParseFolder(uploadsPath)
	fmt.Println(listOfFiles)
	resultJson := analyser.AnalyseList(listOfFiles)
	fmt.Println(resultJson)
	// if err != nil {
	// 	log.Fatal("Start analyse request err:", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Analysis failed"})
	// 	return
	// }

	c.JSON(http.StatusOK, resultJson)
}

// Route /analyse
func Analyse(c *gin.Context) {
	userName := c.PostForm("userName")
	if userName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username missing"})
		return
	}

	StartAnalyse(userName, c)
}
