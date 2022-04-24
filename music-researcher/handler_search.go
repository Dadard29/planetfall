package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/zmb3/spotify/v2"
)

func (s *musicResearcherService) listArtistsFromTrack(
	ctx context.Context,
	track spotify.FullTrack,
	artistBufferList []spotify.FullArtist,

) ([]spotify.FullArtist, error) {

	out := make([]spotify.FullArtist, 0)
	for _, artist := range track.Artists {
		// check if track artist already in buffer
		inBuffer := false
		for _, artistBuffer := range artistBufferList {
			if artist.ID == artistBuffer.ID {
				out = append(out, artistBuffer)
				inBuffer = true
			}
		}

		// if not, request the artist from API
		if !inBuffer {
			artistBuffer, err := s.spotifyClient.GetArtist(ctx, artist.ID)
			if err != nil {
				return nil, err
			}
			out = append(out, *artistBuffer)
		}
	}

	return out, nil
}

func (s *musicResearcherService) pagesTracksToDTO(ctx context.Context, pages *spotify.FullTrackPage) ([]TrackDTO, error) {
	var trackDtoList = make([]TrackDTO, 0)
	var artistBufferList = make([]spotify.FullArtist, 0)

	for {
		for _, track := range pages.Tracks {
			artistList, err := s.listArtistsFromTrack(ctx, track, artistBufferList)
			if err != nil {
				return nil, err
			}

			trackDto := s.newTrackDTO(track, artistList)
			trackDtoList = append(trackDtoList, trackDto)
		}

		if err := s.spotifyClient.NextPage(ctx, pages); err == spotify.ErrNoMorePages {
			break
		}
	}

	return trackDtoList, nil
}

func (s *musicResearcherService) handlerSpotifySearch(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query().Get("q")
	genreParam := r.URL.Query().Get("gl")
	var genreList []string
	if genreParam != "" {
		genreList = strings.Split(genreParam, ",")
	} else {
		genreList = make([]string, 0)
	}

	ctx := context.Background()
	err := s.setSpotifyClient(ctx)
	if err != nil {
		s.server.InternalServerError(w, r, err, "failed setting up connection with Spotify")
		return
	}

	genreFilters := ""
	for _, genre := range genreList {
		genreFilters = fmt.Sprintf("%s genre:%s", genreFilters, genre)
	}
	queryWithFilters := fmt.Sprintf("%s %s", queryParam, genreFilters)

	log.Printf("requesting Spotify with query: %s", queryWithFilters)
	searchResult, err := s.spotifyClient.Search(ctx, queryWithFilters, spotify.SearchTypeTrack)
	if err != nil {
		s.server.InternalServerError(w, r, err, "failed interacting with Spotify")
		return
	}

	pages := searchResult.Tracks
	trackDtoList, err := s.pagesTracksToDTO(ctx, pages)
	if err != nil {
		s.server.InternalServerError(w, r, err, "failed to extract track pages")
	}

	s.server.JsonResponse(w, r, trackDtoList)
}
