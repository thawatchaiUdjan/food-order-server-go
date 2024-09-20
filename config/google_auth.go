package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func LoadGoogleAuth() *oauth2.Config {
	config := LoadConfig()
	googleAuth := &oauth2.Config{
		ClientID:     config.Google.ClientID,
		ClientSecret: config.Google.SecretID,
		RedirectURL:  config.Google.RedirectURL,
		Endpoint:     google.Endpoint,
	}
	return googleAuth
}
