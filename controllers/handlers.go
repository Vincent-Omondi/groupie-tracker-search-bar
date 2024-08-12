// controllers/handler.go
package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"learn.zone01kisumu.ke/git/johnodhiambo0/groupie-tracker/api"
)

// GetArtistsHandler handles the /artists route
func GetArtistsHandler(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists()
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(artists); err != nil {
		http.Error(w, "Failed to encode locations", http.StatusInternalServerError)
	}
}

// GetLocationsHandler handles the /locations route
func GetLocationsHandler(w http.ResponseWriter, r *http.Request) {
	locations, err := api.GetLocations()
	if err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		http.Error(w, "Failed to encode locations", http.StatusInternalServerError)
	}
}

// GetDatesHandler handles the /dates route
func GetDatesHandler(w http.ResponseWriter, r *http.Request) {
	dates, err := api.GetDates()
	if err != nil {
		http.Error(w, "Failed to fetch dates", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dates); err != nil {
		http.Error(w, "Failed to encode locations", http.StatusInternalServerError)
	}
}

// GetRelationsHandler handles the /relations route
func GetRelationsHandler(w http.ResponseWriter, r *http.Request) {
	relations, err := api.GetRelations()
	if err != nil {
		http.Error(w, "Failed to fetch relations", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relations); err != nil {
		http.Error(w, "Failed to encode locations", http.StatusInternalServerError)
	}
}

func GetArtistByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get artist ID from the URL path
	idStr := r.URL.Path[len("/artists/"):]
	artistID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artist, relation, err := api.GetArtistByID(artistID)
	if err != nil {
		http.Error(w, "Failed to fetch artist or relation", http.StatusInternalServerError)
		return
	}

	// Create a response combining artist and relation data
	response := struct {
		Artist   *api.Artist   `json:"artist"`
		Relation *api.Relation `json:"relation"`
	}{
		Artist:   artist,
		Relation: relation,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
