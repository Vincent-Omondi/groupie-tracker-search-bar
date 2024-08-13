package controllers

import (
	"net/http"
)

// controllers/routes.go

func RegisterRoutes() {
	http.HandleFunc("/artists", ServeArtists)
	http.HandleFunc("/artists/", GetArtistByIDHandler)
	http.HandleFunc("/locations", GetLocationsHandler)
	http.HandleFunc("/dates", GetDatesHandler)
	http.HandleFunc("/relations", GetRelationsHandler)
}
