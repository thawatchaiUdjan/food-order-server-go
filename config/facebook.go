package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

func LoadFacebook() *oauth2.Config {
	config := LoadConfig()
	facebook := &oauth2.Config{
		ClientID:     config.Facebook.ClientID,
		ClientSecret: config.Facebook.SecretID,
		RedirectURL:  config.Google.RedirectURL,
		Endpoint:     facebook.Endpoint,
	}
	return facebook
}
