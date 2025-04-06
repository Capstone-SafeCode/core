package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver/v3"
)

// Route /upload
func UploadFile(c *gin.Context) {
	userName := c.PostForm("userName")
	if userName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username missing"})
		return
	}

	file, header, err := c.Request.FormFile("codeFile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	uploadDir := "uploads/" + userName
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
}
