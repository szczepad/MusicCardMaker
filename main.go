package main

import (
	"fmt"
	"log/slog"

	"github.com/szczepad/MusicCardMaker/config"
	"github.com/szczepad/MusicCardMaker/spotify"
)

const (
	spotifyAuthURL = "https://accounts.spotify.com"
	spotifyApiURL  = "https://api.spotify.com"
)

func main() {
	config := config.CreateConfig()

	client := spotify.NewSpotifyClient(
		spotifyAuthURL,
		spotifyApiURL,
		config.Spotify.ClientID,
		config.Spotify.ClientSecret,
	)

	token, err := client.Authenticate()
	if err != nil {
		slog.Error("Could not authenticate to Spotify", "Error", err)
	}
	fmt.Println(token)
}
