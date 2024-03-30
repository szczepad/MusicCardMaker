package spotify

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

const tokenPath = "/api/token"

type SpotifyClient struct {
	authURL      string
	apiURL       string
	clientID     string
	clientSecret string
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   uint   `json:"expires_in"`
}

type TrackResponse struct {
	Tracks TrackItems `json:"tracks"`
}

type TrackItems struct {
	Items []Item `json:"items"`
}

type Item struct {
	Track TrackInfo `json:"track"`
}

type TrackInfo struct {
	Album        AlbumInfo    `json:"album"`
	Artists      []Artist     `json:"artists"`
	ExternalURLs ExternalURLs `json:"external_urls"`
	Name         string       `json:"name"`
}

type AlbumInfo struct {
	ReleaseDate string `json:"release_date"`
}

type Artist struct {
	Name string `json:"name"`
}

type ExternalURLs struct {
	Spotify string `json:"spotify"`
}

type Track struct {
	Artist      string
	Name        string
	Url         string
	ReleaseYear string
}

func NewSpotifyClient(authURL, apiURL, clientID, clientSecret string) SpotifyClient {
	return SpotifyClient{authURL, apiURL, clientID, clientSecret}
}

func (c *SpotifyClient) Authenticate() (string, error) {
	client := &http.Client{}

	data := url.Values{
		"grant_type": {"client_credentials"},
	}
	req, err := http.NewRequest(
		http.MethodPost,
		c.authURL+tokenPath,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		slog.Error("Could not create Request.", "Error", err)
		return "", err
	}
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

func (c *SpotifyClient) GetTracksFromPlaylist(
	accessToken, playlistID string,
) ([]Track, error) {
	result := []Track{}

	client := &http.Client{}

	req, err := http.NewRequest(
		http.MethodGet,
		c.apiURL+"/v1/playlists/"+playlistID+"/tracks",
		nil,
	)
	if err != nil {
		slog.Error("Could not create Request to get Tracks from Playlist", "Error", err)
		return result, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	res, err := client.Do(req)
	if err != nil {
		slog.Error("Could not perform Request to get Tracks from Playlist", "Error", err)
		return result, err
	}
	if res.StatusCode != 200 {
		slog.Error("Got unexpected StatusCode", "Status", res.StatusCode)
		return result, errors.New("InvalidStatus")
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Could not read response Body", "Error", err)
		return result, err
	}

	var trackResponse TrackResponse
	if err := json.Unmarshal(body, &trackResponse); err != nil {
		log.Fatal(err)
	}
	for _, item := range trackResponse.Tracks.Items {
		releaseYear, _, found := strings.Cut(item.Track.Album.ReleaseDate, "-")
		if !found {
			if len(item.Track.Album.ReleaseDate) != 4 {
				slog.Error(
					"Could not get ReleaseYear for Track",
					"ReleaseDate",
					item.Track.Album.ReleaseDate,
				)
			} else {
				releaseYear = item.Track.Album.ReleaseDate
			}
		}

		track := Track{
			Artist:      item.Track.Artists[0].Name,
			Name:        item.Track.Name,
			Url:         item.Track.ExternalURLs.Spotify,
			ReleaseYear: releaseYear,
		}
		result = append(result, track)

	}
	return result, nil
}
