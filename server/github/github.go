package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type GithubUserData struct {
  ID int `json:"id"`
	Login string `json:"login"`
  AvatarURL string `json:"avatar_url"`
}

func ExchangeCodeWithToken(code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", os.Getenv("GH_CLIENT_ID"))
	data.Set("client_secret", os.Getenv("GH_CLIENT_SECRET"))
	data.Set("code", code)
	req, err := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBufferString(data.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("Error while constructing http request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	httpClient := http.Client{
		Timeout: time.Second * 10,
	}
	response, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error while sending http request: %w", err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error while reading response body: %w", err)
	}
	githubResponse := struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}{}
	err = json.Unmarshal(responseBody, &githubResponse)
	if err != nil {
		return "", fmt.Errorf("Error while parsing github response: %w", err)
	}
	return githubResponse.AccessToken, nil
}

func GetGithubUserData(accessToken string) (*GithubUserData, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Error while constructing http request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", accessToken))
	httpClient := http.Client{Timeout: time.Second * 10}
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while sending http request: %w", err)
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while reading http response body: %w", err)
	}
	data := &GithubUserData{}
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing github response: %w", err)
	}
	return data, nil
}
