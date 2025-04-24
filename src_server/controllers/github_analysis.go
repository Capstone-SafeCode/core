package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"test_capstone/src_server/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v58/github"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

// getUserFromToken récupère l'utilisateur à partir du token JWT
func getUserFromToken(c *gin.Context) (*model.User, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return nil, fmt.Errorf("utilisateur non authentifié")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		return nil, fmt.Errorf("ID utilisateur invalide")
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("utilisateur non trouvé")
	}

	if user.GitHubToken == "" {
		return nil, fmt.Errorf("token GitHub non trouvé. Veuillez vous reconnecter avec GitHub")
	}

	return &user, nil
}

// createGitHubClient crée un client GitHub avec le token de l'utilisateur
func createGitHubClient(token string) (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Vérifier que le token est valide
	_, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("token GitHub invalide: %v", err)
	}

	return client, nil
}

// extractRepoInfo extrait les informations du repository à partir de l'URL
func extractRepoInfo(repoURL string) (string, string, error) {
	if !strings.HasPrefix(repoURL, "https://github.com/") {
		return "", "", fmt.Errorf("URL invalide: doit commencer par https://github.com/")
	}

	repoURL = strings.TrimSuffix(repoURL, ".git")
	repoURL = strings.TrimSuffix(repoURL, "/")

	parts := strings.Split(repoURL, "/")
	if len(parts) < 5 {
		return "", "", fmt.Errorf("URL invalide: doit être au format https://github.com/owner/repo")
	}

	owner := parts[3]
	repo := parts[4]

	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("URL invalide: le propriétaire et le nom du repository ne peuvent pas être vides")
	}

	return owner, repo, nil
}

// downloadRepository télécharge un repository GitHub
func downloadRepository(client *github.Client, owner, repo string, isPrivate bool) (string, error) {
	ctx := context.Background()

	// Vérifier les permissions pour les repositories privés
	if isPrivate {
		permission, _, err := client.Repositories.GetPermissionLevel(ctx, owner, repo, owner)
		if err != nil {
			return "", fmt.Errorf("erreur lors de la vérification des permissions: %v", err)
		}

		if *permission.Permission == "none" {
			return "", fmt.Errorf("vous n'avez pas les permissions nécessaires pour accéder à ce repository")
		}
	}

	// Obtenir le lien de téléchargement
	zipURL, _, err := client.Repositories.GetArchiveLink(ctx, owner, repo, github.Zipball, &github.RepositoryContentGetOptions{
		Ref: "main",
	}, 1)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération du repository: %v", err)
	}

	// Télécharger le fichier
	resp, err := http.Get(zipURL.String())
	if err != nil {
		return "", fmt.Errorf("erreur lors du téléchargement: %v", err)
	}
	defer resp.Body.Close()

	// Créer le fichier zip
	zipPath := filepath.Join("temp", repo+".zip")
	if err := os.MkdirAll("temp", os.ModePerm); err != nil {
		return "", fmt.Errorf("erreur lors de la création du dossier temporaire: %v", err)
	}

	out, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la création du fichier zip: %v", err)
	}
	defer out.Close()

	// Copier le contenu
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la copie du contenu: %v", err)
	}

	return zipPath, nil
}

// sendToUpload envoie le fichier zip à la route /upload
func sendToUpload(c *gin.Context, zipPath string, user *model.User) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Ajouter le nom d'utilisateur
	if err := writer.WriteField("userName", user.Username); err != nil {
		return fmt.Errorf("erreur lors de l'ajout du nom d'utilisateur: %v", err)
	}

	// Ajouter le fichier
	file, err := os.Open(zipPath)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ouverture du fichier zip: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("codeFile", filepath.Base(zipPath))
	if err != nil {
		return fmt.Errorf("erreur lors de la création du formulaire: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("erreur lors de la copie du fichier: %v", err)
	}

	writer.Close()

	// Créer la requête
	req, err := http.NewRequest("POST", "http://localhost:8080/upload", body)
	if err != nil {
		return fmt.Errorf("erreur lors de la création de la requête: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", c.GetHeader("Authorization"))

	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'envoi de la requête: %v", err)
	}
	defer resp.Body.Close()

	// Lire la réponse
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erreur lors de la lecture de la réponse: %v", err)
	}

	// Retourner la réponse
	c.Data(resp.StatusCode, "application/json", respBody)
	return nil
}

// DownloadGitHubRepo est la fonction principale qui orchestre le processus
func DownloadGitHubRepo(c *gin.Context) {
	// Récupérer l'utilisateur
	user, err := getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Récupérer l'URL du repository
	repoURL := c.PostForm("repo_url")
	if repoURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL du repository manquante"})
		return
	}

	// Extraire les informations du repository
	owner, repo, err := extractRepoInfo(repoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Créer le client GitHub
	client, err := createGitHubClient(user.GitHubToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Vérifier si le repository existe et est accessible
	repoInfo, _, err := client.Repositories.Get(context.Background(), owner, repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur lors de la vérification du repository",
			"details": fmt.Sprintf("Owner: %s, Repo: %s, Error: %v, Token: %s...",
				owner, repo, err, user.GitHubToken[:10]+"..."),
		})
		return
	}

	// Télécharger le repository
	zipPath, err := downloadRepository(client, owner, repo, *repoInfo.Private)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer os.Remove(zipPath) // Nettoyer le fichier temporaire

	// Envoyer à /upload
	if err := sendToUpload(c, zipPath, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
