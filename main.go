package main

import (
	"log/slog"

	"github.com/szczepad/MusicCardMaker/config"
	"github.com/szczepad/MusicCardMaker/spotify"
)

const (
	spotifyAuthURL = "https://accounts.spotify.com"
	spotifyApiURL  = "https://api.spotify.com"
)

func main() {
	playlistID := "0tarwRmyLGjw3QlMq4GNhn?si=899e9723d2fb483f" // TODO: Make this configurable via command line
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
	tracks, err := client.GetTracksFromPlaylist(token, playlistID)

	spotify.CreatePDF(tracks)
}
