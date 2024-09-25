package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"learn.zone01kisumu.ke/git/johnodhiambo0/groupie-tracker/api"
)

type TemplateData struct {
	Artists   []api.Artist
	Query     string
	NoResults bool
}

type ArtistDetailData struct {
	Artist   api.Artist
	Location []string
	Date     api.Date
	Relation api.Relation
}

var (
	artistCache       []api.Artist
	locationCache     map[int][]string
	cacheTime         time.Time
	locationCacheTime time.Time
)

const cacheDuration = 10 * time.Minute

// ErrorHandler handles error responses and templates
func ErrorHandler(w http.ResponseWriter, message string, statusCode int, logError, showStatusCode bool) {
	w.WriteHeader(statusCode)

	data := struct {
		StatusCode int
		ErrMsg     string
	}{
		StatusCode: statusCode,
		ErrMsg:     message,
	}

	// If showStatusCode is false, set StatusCode to 0 to avoid displaying it
	if !showStatusCode {
		data.StatusCode = 0
	}

	tmpl, err := template.ParseFiles("templates/error.html")
	if err != nil {
		if logError {
			log.Println("Error parsing error template:", err)
		}
		http.Error(w, "Oops! Something went wrong on our end. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		if logError {
			log.Println("Error executing error template:", err)
		}
		return
	}

	_, err = buf.WriteTo(w)
	if err != nil && logError {
		log.Println("Error writing response:", err)
	}
}

// ServeArtists handles the /artists route
func ServeArtists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	// Check if cache is valid
	if time.Since(cacheTime) > cacheDuration || artistCache == nil {
		artists, err := api.GetArtists()
		if err != nil {
			log.Printf("Error getting artists: %v", err)
			ErrorHandler(w, "Oops!\n We ran into an issue while fetching Artists,\n Please try again later.", http.StatusInternalServerError, false, false)
			return
		}
		artistCache = artists
		cacheTime = time.Now()
	}

	// check if location cache is valid
	if time.Since(locationCacheTime) > cacheDuration || locationCache == nil {
		locationCache = make(map[int][]string)
		for _, artist := range artistCache {
			locations, err := api.GetLocationsByArtistID(artist.ID)
			if err != nil {
				log.Printf("Error getting locations for artist %d: %v", artist.ID, err)
				continue
			}
			locationCache[artist.ID] = locations
		}
		locationCacheTime = time.Now()
	}

	filteredArtists := filterArtists(artistCache, query)

	// Check if no results were found and query is not empty
	if len(filteredArtists) == 0 && query != "" {
		ErrorHandler(w, "No Result Found for this search.", http.StatusNotFound, false, false)
		return
	}

	data := TemplateData{
		Artists:   filteredArtists,
		Query:     query,
		NoResults: len(filteredArtists) == 0 && query != "",
	}

	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		ErrorHandler(w, "An unexpected error occurred. Please try again later.", http.StatusInternalServerError, true, true)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, "We encountered an issue while rendering the page. Please try again later.", http.StatusInternalServerError, true, true)
		return
	}
}

// filterArtists filters the list of artists based on the search query
func filterArtists(artists []api.Artist, query string) []api.Artist {
	if query == "" {
		return artists
	}

	query = strings.ToLower(query)
	var result []api.Artist

	for _, a := range artists {
		// Artist/band name matches
		if strings.Contains(strings.ToLower(a.Name), query) {
			result = append(result, a)
			continue
		}

		// Members
		for _, member := range a.Members {
			if strings.Contains(strings.ToLower(member), query) {
				result = append(result, a)
				break
			}
		}

		// Locations
		locations := strings.Split(a.Locations, ",")
		for _, location := range locations {
			if strings.Contains(strings.ToLower(strings.TrimSpace(location)), query) {
				result = append(result, a)
				break
			}
		}

		// First album dates
		if strings.Contains(strings.ToLower(a.FirstAlbum), query) {
			result = append(result, a)
			continue
		}

		// creation dates
		if strings.Contains(strconv.Itoa(a.CreationDate), query) {
			result = append(result, a)
			continue
		}
	}
	return result
}

func GetSearchSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == " " {
		json.NewEncoder(w).Encode([]string{})
		return
	}
	suggestions := []string{}
	artists, _ := api.GetArtists()

	for _, artist := range artists {
		// Artist/band name
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			suggestions = append(suggestions, fmt.Sprintf("%s - artist/band", artist.Name))
		}

		// Members
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), strings.ToLower(query)) {
				suggestions = append(suggestions, fmt.Sprintf("%s - member", member))
			}
		}

		// Locations
		locations, _ := api.GetLocationsByArtistID(artist.ID)
		// locations := strings.Split(artistLocations,",")
		for _, location := range locations {
			// trimmedLocations := strings.TrimSpace(location)
			if strings.Contains(strings.ToLower(strings.TrimSpace(location)), strings.ToLower(query)) {
				suggestions = append(suggestions, fmt.Sprintf("%s - location", location))
			}
		}

		// First album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(query)) {
			suggestions = append(suggestions, fmt.Sprintf("%s - first album date", artist.FirstAlbum))
		}

		// Creation date
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			suggestions = append(suggestions, fmt.Sprintf("%d - creation date", artist.CreationDate))
		}
	}

	json.NewEncoder(w).Encode(suggestions)

}

func ServeArtistDetails(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/artist/" {
		ErrorHandler(w, "oops! page not found", http.StatusNotFound, true, true)
		return
	}

	idStr := r.URL.Path[len("/artist/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorHandler(w, "Invalid artist ID", http.StatusBadRequest, true, true)
		return
	}

	artist, _, date, relation, err := api.GetArtistByID(id)
	if err != nil {
		log.Printf("Error retrieving artist by ID %v: %s", id, err)
		ErrorHandler(w, "Ooops!\n We ran into an issue while fetching Artists,\n Please try again later.", http.StatusInternalServerError, false, false)
		return
	}

	locations, err := api.GetLocationsByArtistID(id)
	if err != nil {
		log.Printf("Error retrieving locations for artist ID %v: %s", id, err)
		ErrorHandler(w, "Ooops! Could not retrieve locations.", http.StatusInternalServerError, false, false)
		return
	}

	data := ArtistDetailData{
		Artist:   *artist,
		Location: locations,
		Date:     *date,
		Relation: *relation,
	}

	tmpl, err := template.ParseFiles("templates/artist_details.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		ErrorHandler(w, "Unable to load artist details at this time. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, "Error rendering artist details. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

// GetArtistsHandler handles the /artists route
func GetArtistsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	artists, err := api.GetArtists()
	if err != nil {
		log.Printf("Error fetching artists: %v", err)
		ErrorHandler(w, "Unable to retrieve artist information at this time. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}

	filteredArtists := filterArtists(artists, query)

	if len(filteredArtists) == 0 {
		ErrorHandler(w, "No artists found matching the search term.", http.StatusNotFound, false, false)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(filteredArtists); err != nil {
		log.Printf("Error encoding artists data to JSON: %v", err)
		ErrorHandler(w, "An error occurred while processing the artist data. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

// GetLocationsHandler handles the /locations route
func GetLocationsHandler(w http.ResponseWriter, r *http.Request) {
	locations, err := api.GetLocations()
	if err != nil {
		log.Printf("Error fetching locations: %v", err)
		ErrorHandler(w, "Unable to retrieve locations at this time. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		log.Printf("Error encoding locations data to JSON: %v", err)
		ErrorHandler(w, "An error occurred while processing location data. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

// GetDatesHandler handles the /dates route
func GetDatesHandler(w http.ResponseWriter, r *http.Request) {
	dates, err := api.GetDates()
	if err != nil {
		log.Printf("Error fetching dates: %v", err)
		ErrorHandler(w, "Unable to retrieve dates at this time. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dates); err != nil {
		log.Printf("Error encoding dates data to JSON: %v", err)
		ErrorHandler(w, "An error occurred while processing date information. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

// GetRelationsHandler handles the /relations route
func GetRelationsHandler(w http.ResponseWriter, r *http.Request) {
	relations, err := api.GetRelations()
	if err != nil {
		log.Printf("Error fetching relations: %v", err)
		ErrorHandler(w, "Unable to retrieve relations at this time. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(relations); err != nil {
		log.Printf("Error encoding relations data to JSON: %v", err)
		ErrorHandler(w, "An error occurred while processing relation data. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

func GetArtistByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get artist ID from the URL path
	idStr := r.URL.Path[len("/artists/"):]
	artistID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid artist ID: %v", err)
		ErrorHandler(w, "Invalid artist ID provided. Please check and try again.", http.StatusBadRequest, true, true)
		return
	}

	artist, location, date, relation, err := api.GetArtistByID(artistID)
	if err != nil {
		// Check if the artist was not found
		if err.Error() == "artist not found" {
			ErrorHandler(w, "Artist not found. Please check the ID and try again.", http.StatusNotFound, true, true)
		} else {
			log.Printf("Error fetching artist details with ID %d: %v", artistID, err)
			ErrorHandler(w, "Unable to retrieve artist details at this time. Please try again later.", http.StatusInternalServerError, false, false)
		}
		return
	}

	// Create a response combining artist, location, date, and relation data
	response := struct {
		Artist   *api.Artist   `json:"artist"`
		Location *api.Location `json:"location"`
		Date     *api.Date     `json:"date"`
		Relation *api.Relation `json:"relation"`
	}{
		Artist:   artist,
		Location: location,
		Date:     date,
		Relation: relation,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response for artist ID %d: %v", artistID, err)
		ErrorHandler(w, "An error occurred while processing the response. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

// Serve About Page
func AboutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorHandler(w, "method not allowed", http.StatusMethodNotAllowed, true, true)
		return
	}
	tmpl, err := template.ParseFiles("templates/about.html")
	if err != nil {
		ErrorHandler(w, "file not found", http.StatusNotFound, true, true)
		log.Printf("Error parsing index.html: %v\n", err)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		ErrorHandler(w, "Something unexpected occured", http.StatusInternalServerError, false, false)
		log.Printf("Error executing template: %v\n", err)
		return
	}
}
