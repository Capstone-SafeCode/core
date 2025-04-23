package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"test_capstone/src_server/config"

	"github.com/gin-gonic/gin"
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

	// Ici, vous pouvez stocker les informations de l'utilisateur dans votre base de données
	// et créer une session pour l'utilisateur

	c.JSON(http.StatusOK, gin.H{
		"message": "Authentification réussie",
		"user":    user,
	})
}
