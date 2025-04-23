package main

import (
	"log"
	"os"
	"path/filepath"
	"test_capstone/src_server/controllers"
	"test_capstone/src_server/database"
	"test_capstone/src_server/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Trouver le chemin absolu du répertoire de travail
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Erreur lors de la récupération du répertoire de travail")
	}

	// Charger les variables d'environnement
	envPath := filepath.Join(wd, ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Attention: Impossible de charger le fichier .env: %v", err)
	}

	database.ConnectDB()

	controllers.InitCollections()

	router := routes.SetupRouter()
	router.Run(":8080")
}
