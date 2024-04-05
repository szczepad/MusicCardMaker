package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/szczepad/MusicCardMaker/config"
	"github.com/szczepad/MusicCardMaker/spotify"
)

const (
	spotifyAuthURL = "https://accounts.spotify.com"
	spotifyApiURL  = "https://api.spotify.com"
)

func main() {
	//:= "0tarwRmyLGjw3QlMq4GNhn?si=899e9723d2fb483f" // TODO: Make this configurable via command line
	config := config.CreateConfig()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, " %s [options] <Spotify playlist or album link>\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(
			os.Stderr,
			" %s https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M\n",
			os.Args[0],
		)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		slog.Error("Error: A Spotify playlist or album link is required")
		flag.Usage()
		os.Exit(1)
	}

	playlistID := flag.Arg(0)

	client := spotify.NewSpotifyClient(
		spotifyAuthURL,
		spotifyApiURL,
		config.Spotify.ClientID,
		config.Spotify.ClientSecret,
	)

	token, err := client.Authenticate()
	if err != nil {
		slog.Error("Could not authenticate to Spotify", "Error", err)
		os.Exit(1)
	}
	tracks, err := client.GetTracksFromPlaylist(token, playlistID)
	if err != nil {
		slog.Error("Could not get Tracks from Playlist", "Error", err)
		os.Exit(1)
	}

	err = spotify.CreatePDF(tracks)
	if err != nil {
		slog.Error("Could not create PDF for Tracks", "Error", err)
		os.Exit(1)
	}
}
