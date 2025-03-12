package controllers

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"test_capstone/src_server/database"
	"test_capstone/src_server/model"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var userCollection *mongo.Collection

func InitCollections() {
	userCollection = database.DB.Collection("users")
}

// Route /register
func CreateUser(c *gin.Context) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Route /login
func LogUser(c *gin.Context) {
	var user model.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result model.User
	err := userCollection.FindOne(ctx, bson.M{"Username": user.Username}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	if result.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Route /users
func GetUsers(c *gin.Context) {
	var users []model.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user model.User
		cursor.Decode(&user)
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}
