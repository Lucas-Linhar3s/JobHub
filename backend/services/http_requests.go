package services

import (
	"strconv"
	"time"

	"github.com/Lucas-Linhar3s/JobHub/pkg/log"
	"github.com/Lucas-Linhar3s/JobHub/utils"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func GetUserDataRequest(provider *string, token *oauth2.Token, logger *log.Logger) (interface{}, error) {
	switch *provider {
	case "Google":
		return GetDataUserGoogle(token, logger)
	case "Github":
		return GetDataUserGithub(token, logger)
	case "Linkedin":
		return GetDataUserLinkedin(token, logger)
	}
	return nil, nil
}

func GetDataUserGithub(token *oauth2.Token, logger *log.Logger) (*GitHubRes, error) {
	var (
		data        = new(GitHubRes)
		urlUserData = "https://api.github.com/user"
	)

	headers := map[string][]string{
		"Authorization": {"Bearer " + token.AccessToken},
		"Content-Type":  {"application/json"},
		"Accept":        {"application/json"},
	}

	httpClient := utils.NewHTTPClient(urlUserData)
	reqGET, err := httpClient.BuildGetRequest("Github Get User Data", "", headers, nil)
	if err != nil {
		return nil, err
	}

	// var errorData map[string]interface{}
	logger.Info(
		"Requisição GetDataUserGithub Iniciada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserData))

	err = httpClient.PerformRequest(reqGET, data, nil)
	if err != nil {
		return nil, err
	}

	emails, err := getUserEmailsGithub(token, logger)
	if err != nil {
		return nil, err
	}

	for _, email := range emails {
		if email.Primary != nil && *email.Primary {
			data.Email = email.Email
			break
		}
	}
	data.OauthId = utils.GetStringPointer(strconv.Itoa(*data.Id))

	logger.Info(
		"Requisição GetDataUserGithub Finalizada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserData))

	return data, nil
}

func getUserEmailsGithub(token *oauth2.Token, logger *log.Logger) ([]*GithubEmails, error) {
	var (
		data          []*GithubEmails
		urlUserEmails = "https://api.github.com/user/emails"
	)

	headers := map[string][]string{
		"Authorization": {"Bearer " + token.AccessToken},
		"Content-Type":  {"application/json"},
	}

	httpClient := utils.NewHTTPClient(urlUserEmails)
	reqGET, err := httpClient.BuildGetRequest("Github Get User Data", "", headers, nil)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"Requisição getUserEmailsGithub Iniciada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserEmails))

	err = httpClient.PerformRequest(reqGET, &data, nil)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"Requisição getUserEmailsGithub Finalizada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserEmails))

	return data, nil
}

func GetDataUserGoogle(token *oauth2.Token, logger *log.Logger) (*GoogleRes, error) {
	var (
		data        = new(GoogleRes)
		urlUserData = "https://www.googleapis.com/oauth2/v2/userinfo"
	)

	headers := map[string][]string{
		"Authorization": {"Bearer " + token.AccessToken},
		"Content-Type":  {"application/json"},
	}

	httpClient := utils.NewHTTPClient(urlUserData)
	reqGET, err := httpClient.BuildGetRequest("Google Get User Data", "", headers, nil)
	if err != nil {
		return nil, err
	}

	// var errorData map[string]interface{}
	logger.Info(
		"Requisição GetDataUserGoogle Iniciada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserData))

	err = httpClient.PerformRequest(reqGET, data, nil)
	if err != nil {
		return nil, err
	}
	data.OauthId = data.Id

	logger.Info(
		"Requisição GetDataUserGithub Finalizada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserData))

	return data, nil
}

func GetDataUserLinkedin(token *oauth2.Token, logger *log.Logger) (interface{}, error) {
	var (
		data        interface{}
		urlUserData = "https://api.linkedin.com/v2/me"
	)

	headers := map[string][]string{
		"Authorization": {"Bearer " + token.AccessToken},
		"Content-Type":  {"application/json"},
	}

	httpClient := utils.NewHTTPClient(urlUserData)
	reqGET, err := httpClient.BuildGetRequest("Google Get User Data", "", headers, nil)
	if err != nil {
		return nil, err
	}

	// var errorData map[string]interface{}
	logger.Info(
		"Requisição GetDataUserGoogle Iniciada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserData))

	err = httpClient.PerformRequest(reqGET, data, nil)
	if err != nil {
		return nil, err
	}
	// data.OauthId = data.Id

	logger.Info(
		"Requisição GetDataUserGithub Finalizada",
		zap.Time("time", time.Now()),
		zap.String("url", urlUserData))

	return data, nil
}
