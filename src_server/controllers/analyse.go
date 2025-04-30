package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	analyser "test_capstone/src_analyser"
	parser "test_capstone/src_parser"
	"test_capstone/src_server/database"
	"test_capstone/src_server/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// Récupérer l'ID de l'utilisateur depuis le token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convertir l'ID en ObjectID
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Créer une nouvelle analyse
	analysis := model.Analysis{
		UserID:    objID,
		Timestamp: primitive.NewDateTimeFromTime(time.Now()),
		Results:   resultJson,
	}

	// Sauvegarder dans MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = database.DB.Collection("analyses").InsertOne(ctx, analysis)
	if err != nil {
		log.Printf("Error saving analysis: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save analysis history"})
		return
	}

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

// GetAnalysisHistory récupère l'historique des analyses d'un utilisateur
func GetAnalysisHistory(c *gin.Context) {
	// Récupérer l'ID de l'utilisateur depuis le token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Convertir l'ID en ObjectID
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Récupérer l'historique des analyses
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := database.DB.Collection("analyses").Find(ctx, bson.M{"user_id": objID})
	if err != nil {
		log.Printf("Error retrieving analysis history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve analysis history"})
		return
	}
	defer cursor.Close(ctx)

	var analyses []model.Analysis
	if err = cursor.All(ctx, &analyses); err != nil {
		log.Printf("Error decoding analysis history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode analysis history"})
		return
	}

	c.JSON(http.StatusOK, analyses)
}
