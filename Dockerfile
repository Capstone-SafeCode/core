# Utilisation de l'image officielle Go
FROM golang:1.24.1

# Définition du répertoire de travail
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copier le code source
COPY . .

# Exposer le port 8080
EXPOSE 8080

# Lancer l'application avec `go run`
CMD ["go", "run", "./src_server/main.go"]