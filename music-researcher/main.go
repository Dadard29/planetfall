package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Dadard29/planetfall/planetfall"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type service struct {
	server *planetfall.Server

	spotifyClient *spotify.Client
	spotifyToken  *oauth2.Token
}

var projectID = os.Getenv("PROJECT_ID")
var serviceName = os.Getenv("SERVICE")

const (
	spotifyClientID     = "SPOTIFY_CLIENT_ID"
	spotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
)

func main() {
	svc := service{
		spotifyClient: nil,
		spotifyToken:  nil,
		server:        nil,
	}

	serv, err := planetfall.NewServer(projectID, serviceName, []planetfall.Route{
		planetfall.Route{
			Endpoint: "spotify-search",
			Handler:  svc.handlerSpotifySearch,
			Methods:  []string{http.MethodGet},
		},
	})
	if err != nil {
		log.Panicf("failed to create server: %v", err)
	}

	svc.server = serv
	defer svc.server.close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	serv.listen(":" + port)
}
