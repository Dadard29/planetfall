package main

import (
	"log"
	"os"

	"github.com/Dadard29/planetfall/planetfall"
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type Service struct {
	planetfall.server

	spotifyClient *spotify.Client
	spotifyToken  *oauth2.Token
}

var projectID = os.Getenv("PROJECT_ID")
var service = os.Getenv("SERVICE")

const (
	spotifyClientID     = "SPOTIFY_CLIENT_ID"
	spotifyClientSecret = "SPOTIFY_CLIENT_SECRET"
)

func main() {

	serv, err := planetfall.NewServer()
	if err != nil {
		log.Panicf("failed to create server: %v", err)
	}

	defer serv.close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	serv.listen(":" + port)
}
