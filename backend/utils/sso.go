package utils

import (
	"github.com/Lucas-Linhar3s/JobHub/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

func GetOauth2Config(provider string, config *config.Config) *oauth2.Config {
	var ouathConfig *oauth2.Config
	switch provider {
	case "Google":
		ouathConfig = getSsoGoogle(config)
	case "Github":
		ouathConfig = getSsoGithub(config)
	case "Linkedin":
		ouathConfig = getSsoLinkedin(config)
	}
	return ouathConfig
}

func getSsoGoogle(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.Security.Oauth2.Google.ClientId,
		ClientSecret: config.Security.Oauth2.Google.ClientSecret,
		RedirectURL:  config.Security.Oauth2.Google.RedirectUrl,
		Endpoint:     google.Endpoint,
		Scopes:       config.Security.Oauth2.Google.Scopes,
	}
}

func getSsoGithub(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.Security.Oauth2.Github.ClientId,
		ClientSecret: config.Security.Oauth2.Github.ClientSecret,
		RedirectURL:  config.Security.Oauth2.Github.RedirectUrl,
		Endpoint:     github.Endpoint,
		Scopes:       config.Security.Oauth2.Github.Scopes,
	}
}

func getSsoLinkedin(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.Security.Oauth2.Linkedin.ClientId,
		ClientSecret: config.Security.Oauth2.Linkedin.ClientSecret,
		RedirectURL:  config.Security.Oauth2.Linkedin.RedirectUrl,
		Endpoint:     linkedin.Endpoint,
		Scopes:       config.Security.Oauth2.Linkedin.Scopes,
	}
}
