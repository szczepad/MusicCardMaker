package spotify

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

const tokenPath = "/api/token"

type SpotifyClient struct {
	baseURL      string
	clientID     string
	clientSecret string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

func NewSpotifyClient(baseURL, clientID, clientSecret string) SpotifyClient {
	return SpotifyClient{baseURL, clientID, clientSecret}
}

func (c *SpotifyClient) Authenticate() (string, error) {
	client := &http.Client{}

	data := url.Values{
		"grant_type": {"client_credentials"},
	}
	req, err := http.NewRequest(
		http.MethodPost,
		c.baseURL+tokenPath,
		strings.NewReader(data.Encode()),
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.clientSecret)

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Error while Authenticating", "Error", err)
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return "", errors.New("Unexpected Statuscode. Got:" + res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Join(errors.New("Could not read Body of Authentication request"), err)
	}

	var authResponse AuthResponse
	err = json.Unmarshal(body, &authResponse)
	if err != nil {
		return "", err
	}

	return authResponse.AccessToken, nil
}
