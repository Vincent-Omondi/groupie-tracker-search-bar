package controllers

import (
	"net/http"
)

func RegisterRoutes() {
	http.HandleFunc("/artists", ServeArtists)
	http.HandleFunc("/artists/", GetArtistByIDHandler)
	http.HandleFunc("/locations", GetLocationsHandler)
	http.HandleFunc("/dates", GetDatesHandler)
	http.HandleFunc("/relations", GetRelationsHandler)
	http.HandleFunc("/search-suggestions", GetSearchSuggestionsHandler)

}
