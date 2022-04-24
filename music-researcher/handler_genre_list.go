package main

import (
	"context"
	"net/http"
)

func (s *musicResearcherService) handlerSpotifyGenreList(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	err := s.setSpotifyClient(ctx)
	if err != nil {
		s.server.InternalServerError(w, r, err, "failed setting up connection with Spotify")
		return
	}

	genreList, err := s.spotifyClient.GetAvailableGenreSeeds(ctx)
	if err != nil {
		s.server.InternalServerError(w, r, err, "failed to retrieve genre list")
		return
	}
	s.server.JsonResponse(w, r, genreList)
}
