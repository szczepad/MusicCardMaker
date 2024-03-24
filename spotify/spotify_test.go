package spotify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
						json.NewEncoder(w).Encode(payload)
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

	spotifyClient := NewSpotifyClient(testServer.URL, "testID", "testSecret")
	token, err := spotifyClient.Authenticate()
	if err != nil {
		t.Errorf("Could not authenticate. Error: %v", err)
	}
	if token == "" {
		t.Errorf("Could not authenticate. Got empty token")
	}
}
