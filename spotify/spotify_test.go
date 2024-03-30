package spotify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAuthentication(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/token" && r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				t.Errorf("Could not parse form. Error: %v", err)
			}
			data := r.PostForm
			if data["grant_type"][0] == "client_credentials" {
				user, password, ok := r.BasicAuth()
				if !ok {
					w.WriteHeader(401)
				} else {
					if user == "testID" && password == "testSecret" {
						w.WriteHeader(200)
						payload := AuthResponse{"token", "Bearer", 1234}
						err := json.NewEncoder(w).Encode(payload)
						if err != nil {
							w.WriteHeader(500)
						}
					} else {
						w.WriteHeader(401)
					}
				}
			} else {
				w.WriteHeader(500)
			}
		}
	}))
	defer testServer.Close()

	spotifyClient := NewSpotifyClient(testServer.URL, "", "testID", "testSecret")
	token, err := spotifyClient.Authenticate()
	if err != nil {
		t.Errorf("Could not authenticate. Error: %v", err)
	}
	if token == "" {
		t.Errorf("Could not authenticate. Got empty token")
	}
}

func TestGetTracksFromPlaylist(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/v1/playlists/") && r.Method == "GET" {
			pathSegments := strings.Split(r.URL.Path, "/")

			if len(pathSegments) == 5 && pathSegments[4] == "tracks" {
				playlistID := pathSegments[3]
				if r.Header.Get("Authorization") != "Bearer ValidToken" {
					w.WriteHeader(401)
					return
				}
				file, err := os.ReadFile("../testdata/" + playlistID + ".json")
				if err != nil {
					w.WriteHeader(500)
					return
				}
				_, err = w.Write(file)
				if err != nil {
					w.WriteHeader(500)
					return
				}
				return
			}
		}
	}))
	spotifyClient := NewSpotifyClient("", testServer.URL, "testID", "testSecret")

	tests := []struct {
		name       string
		input      string
		wantTracks []Track
	}{
		{
			name:  "Gets Tracks from an existing Playlist",
			input: "singleTrackPlaylist", wantTracks: []Track{
				{
					Artist:      "Emei",
					Name:        "Irresponsible",
					Url:         "https://open.spotify.com/track/60SugyNV4FdewZfktXfXte",
					ReleaseYear: "2023",
				},
			},
		},
		{
			name:  "Handles cases correctly in which only the releaseYear is provided",
			input: "noReleaseMonthOrDay",
			wantTracks: []Track{
				{
					Artist:      "Aerosmith",
					Name:        "I Don't Want to Miss a Thing",
					Url:         "https://open.spotify.com/intl-de/track/225xvV8r1yKMHErSWivnow?si=b10585f9d2bf4225",
					ReleaseYear: "1998",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			token := "ValidToken"
			gotTracks, err := spotifyClient.GetTracksFromPlaylist(token, tc.input)
			if err != nil {
				t.Errorf("Could not get Tracks from Playlist. Error: %v", err)
			}
			if len(gotTracks) != len(tc.wantTracks) {
				t.Errorf(
					"Did not get correct number of Tracks. Got: %d, Want: %d",
					len(gotTracks),
					len(tc.wantTracks),
				)
			}
			for i, track := range gotTracks {
				if track != tc.wantTracks[i] {
					t.Errorf(
						"Did not get expected Track. Got: %v, Want: %v",
						track,
						tc.wantTracks[i],
					)
				}
			}
		})
	}

	t.Run("Returns an Error if the User is unauthenticated", func(t *testing.T) {
		token := "InvalidToken"

		_, err := spotifyClient.GetTracksFromPlaylist(token, "singleTrackPlaylist")
		if err == nil {
			t.Errorf("Got no error although one was expected.")
		}
	})
}
