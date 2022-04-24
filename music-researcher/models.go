package main

import "github.com/zmb3/spotify/v2"

const spotifyUrlKey = "spotify"

type ArtistDTO struct {
	Name       string `json:"name"`
	SpotifyUrl string `json:"spotify_url"`

	Genres   []string `json:"genres"`
	ImageUrl string   `json:"image_url"`
}

type AlbumDTO struct {
	Name       string `json:"name"`
	SpotifyUrl string `json:"spotify_url"`

	ReleaseDate string `json:"release_date"`
	ImageUrl    string `json:"image_url"`
}

type TrackDTO struct {
	Name       string `json:"name"`
	SpotifyUrl string `json:"spotify_url"`
	DurationMs int    `json:"duration_ms"`
	PreviewURL string `json:"preview_url"`
	Popularity int    `json:"popularity"`

	Album   AlbumDTO    `json:"album"`
	Artists []ArtistDTO `json:"artists"`
}

const (
	typeUnknown = "unknown"
	typeAlbum   = "album"
	typeArtist  = "artist"
)

// fixme: https://github.com/Dadard29/planetfall/issues/4
var defaultImageUrl = map[string]string{
	typeUnknown: "",
	typeAlbum:   "",
	typeArtist:  "",
}

func (s *musicResearcherService) getImageUrl(images []spotify.Image, itemType string) string {
	if len(images) == 0 {
		if url, check := defaultImageUrl[itemType]; !check {
			return defaultImageUrl[typeUnknown]
		} else {
			return url
		}
	}

	return images[0].URL
}

func (s *musicResearcherService) newTrackDTO(track spotify.FullTrack, artistList []spotify.FullArtist) TrackDTO {
	albumDto := AlbumDTO{
		Name:        track.Album.Name,
		ReleaseDate: track.Album.ReleaseDate,
		SpotifyUrl:  track.Album.ExternalURLs[spotifyUrlKey],
		ImageUrl:    s.getImageUrl(track.Album.Images, typeAlbum),
	}

	artistDtoList := make([]ArtistDTO, 0)
	for _, artist := range artistList {
		artistDtoList = append(artistDtoList, ArtistDTO{
			Name:       artist.Name,
			SpotifyUrl: artist.ExternalURLs[spotifyUrlKey],
			Genres:     artist.Genres,
			ImageUrl:   s.getImageUrl(artist.Images, typeArtist),
		})
	}

	trackDto := TrackDTO{
		Name:       track.Name,
		SpotifyUrl: track.ExternalURLs[spotifyUrlKey],
		DurationMs: track.Duration,
		PreviewURL: track.PreviewURL,
		Popularity: track.Popularity,

		Album:   albumDto,
		Artists: artistDtoList,
	}

	return trackDto
}
