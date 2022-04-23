package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zmb3/spotify/v2"
)

func (s *service) handlerSpotifySearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	ctx := context.Background()
	err := s.setSpotifyClient(ctx)
	if err != nil {
		s.server.internalServerError(w, r, err, "failed setting up connection with Spotify")
		return
	}

	searchResult, err := s.spotifyClient.Search(ctx, query, spotify.SearchTypeTrack)
	if err != nil {
		s.server.internalServerError(w, r, err, "failed interacting with Spotify")
		return
	}

	w.WriteHeader(http.StatusOK)
	out := fmt.Sprintf("Got %d tracks from Spotify", len(searchResult.Tracks.Tracks))
	fmt.Fprint(w, out)
}
