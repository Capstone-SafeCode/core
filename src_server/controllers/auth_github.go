package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"test_capstone/src_server/config"
	"test_capstone/src_server/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
)

// GitHubLogin redirige vers la page d'authentification GitHub
func GitHubLogin(c *gin.Context) {
	url := config.OAuth2Config.AuthCodeURL(config.OAuthStateString, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GitHubCallback gère le callback après l'authentification GitHub
func GitHubCallback(c *gin.Context) {
	// Vérifier le state pour la protection CSRF
	state := c.Query("state")
	if state != config.OAuthStateString {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "État invalide"})
		return
	}

	code := c.Query("code")
	token, err := config.OAuth2Config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de l'échange du code"})
		return
	}

	// Utiliser le token pour récupérer les informations de l'utilisateur
	client := config.OAuth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des informations utilisateur"})
		return
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du décodage des informations utilisateur"})
		return
	}

	// Créer un nouvel utilisateur avec les informations GitHub
	githubUser := model.User{
		Username:    user["login"].(string),
		Password:    "",                // Pas de mot de passe pour les utilisateurs GitHub
		GitHubToken: token.AccessToken, // Stocker le token GitHub
	}

	// Vérifier si l'utilisateur existe déjà
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingUser model.User
	err = userCollection.FindOne(ctx, bson.M{"username": githubUser.Username}).Decode(&existingUser)
	if err != nil {
		// L'utilisateur n'existe pas, le créer
		_, err = userCollection.InsertOne(ctx, githubUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création de l'utilisateur"})
			return
		}
		existingUser = githubUser
	} else {
		// Mettre à jour le token GitHub
		_, err = userCollection.UpdateOne(ctx,
			bson.M{"username": githubUser.Username},
			bson.M{"$set": bson.M{"github_token": token.AccessToken}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la mise à jour du token"})
			return
		}
		existingUser.GitHubToken = token.AccessToken
	}

	// Générer le token JWT
	jwtToken, err := generateJWT(existingUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentification GitHub réussie",
		"token":   jwtToken,
		"user":    existingUser,
	})
}
