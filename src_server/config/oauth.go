package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// OAuth2Config contient la configuration OAuth
var OAuth2Config = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/auth/github/callback",
	Scopes:       []string{"user:email"},
	Endpoint:     github.Endpoint,
}

// Chaîne de protection CSRF (state)
const OAuthStateString = "randomStateString"
