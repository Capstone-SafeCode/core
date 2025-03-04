package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
)

func main() {
	router := gin.Default()

	// Middleware CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	//ROUTE POUR L'UPLOAD
	router.POST("/upload", func(c *gin.Context) {
		file, header, err := c.Request.FormFile("codeFile")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
			return
		}
		defer file.Close()

		uploadDir := "uploads"
		os.RemoveAll(uploadDir)
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
			return
		}
		uploadedFilePath := filepath.Join(uploadDir, header.Filename)

		out, err := os.Create(uploadedFilePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file content"})
			return
		}

		extractDir := filepath.Join(uploadDir, header.Filename+"_extracted")
		if err := os.MkdirAll(extractDir, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create extraction directory"})
			return
		}

		err = archiver.Unarchive(uploadedFilePath, extractDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to extract zip: %v", err)})
			return
		}

		if err := os.Remove(uploadedFilePath); err != nil {
			fmt.Printf("Warning: failed to delete zip file: %v\n", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "File uploaded and extracted successfully",
			"extractedPath": extractDir,
		})
	})

	// ROUTE POUR L'ANALYSE
	router.POST("/start_analyse", func(c *gin.Context) {
		uploadsName := "uploads/" + getUploadsName() + "/"
		if uploadsName == "uploads//" {
			return
		}

		cmd := exec.Command("./exec.sh", uploadsName)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		if err != nil {
			log.Fatal("Start analyse request err:", err)
		}

		newJson, err := LoadJsonFromResultFile("result.json")
		if err != nil {
			log.Fatal("Error : result.json doesn't exist", err)
		}

		c.JSON(http.StatusOK, newJson)
	})

	router.Run(":8069")
}
