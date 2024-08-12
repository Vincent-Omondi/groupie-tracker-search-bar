//controllers/routes.go

package controllers

import (
	"net/http"
)

// controllers/routes.go

func RegisterRoutes() {
	http.HandleFunc("/artists", GetArtistsHandler)
	http.HandleFunc("/artists/", GetArtistByIDHandler) // New route for artist by ID
	http.HandleFunc("/locations", GetLocationsHandler)
	http.HandleFunc("/dates", GetDatesHandler)
	http.HandleFunc("/relations", GetRelationsHandler)
}
