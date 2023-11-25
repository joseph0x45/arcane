package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

func GetGitHubAccessToken(code string, logger *logrus.Logger) (token string, err error) {
  data := url.Values{}
  data.Set("client_id", os.Getenv("GITHUB_CLIENT_ID"))
  data.Set("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
  data.Set("code", code)
	if err != nil {
		logger.Error(err)
		return
	}
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBufferString(data.Encode()),
	)
  req.Header.Set("Accept", "application/json")
	if err != nil {
		logger.Error(err)
		return
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	githubResponse := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}{}
	err = json.Unmarshal(respBytes, &githubResponse)
	if err != nil {
		logger.Error(err)
		return
	}
	token = githubResponse.AccessToken
	return
}

func GetGithubData(token string, logger *logrus.Logger) (data map[string]interface{}, err error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		logger.Error(err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	err = json.Unmarshal(respBytes, &data)
	if err != nil {
		logger.Error(err)
	}
	return
}
