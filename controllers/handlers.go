package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
	Location api.Location
	Date     api.Date
	Relation api.Relation
}

var (
	artistCache        []api.Artist
	locationCache      []api.Location
	dateCache          []api.Date
	relationCache      []api.Relation
	cacheTime          time.Time
	cacheMutex         sync.RWMutex
	isCacheInitialized bool
)

const cacheDuration = 10 * time.Minute

func initCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if !isCacheInitialized {
		updateCache()
		isCacheInitialized = true
	}
}

func updateCache() {
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		artists, err := api.GetArtists()
		if err == nil {
			artistCache = artists
		}
	}()

	go func() {
		defer wg.Done()
		locations, err := api.GetLocations()
		if err == nil {
			locationCache = locations
		}
	}()

	go func() {
		defer wg.Done()
		dates, err := api.GetDates()
		if err == nil {
			dateCache = dates
		}
	}()

	go func() {
		defer wg.Done()
		relations, err := api.GetRelations()
		if err == nil {
			relationCache = relations
		}
	}()

	wg.Wait()
	cacheTime = time.Now()
}

func getCachedData() ([]api.Artist, []api.Location, []api.Date, []api.Relation) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if time.Since(cacheTime) > cacheDuration {
		go updateCache()
	}

	return artistCache, locationCache, dateCache, relationCache
}

// ErrorHandler handles error responses and templates
func ErrorHandler(w http.ResponseWriter, message string, statusCode int, logError, showStatusCode bool) {

	if w.Header().Get("Content-Type") != "" {
		// Headers already sent, just log the error and return
		if logError {
			log.Printf("Error occurred after headers were sent: %s", message)
		}
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

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
	if err != nil {
		if logError {
			if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "connection reset by peer") {
				log.Println("Client disconnected before response was fully sent")
			} else {
				log.Println("Error writing response:", err)
			}
		}
		return
	}
}

// ServeArtists handles the /artists route
func ServeArtists(w http.ResponseWriter, r *http.Request) {
	initCache()
	query := r.URL.Query().Get("query")

	artists, _, _, _ := getCachedData()

	for i := range artists {
		// Format the URL with the value of i
		url := fmt.Sprintf("https://groupietrackers.herokuapp.com/api/locations/%d", i+1)

		// Fetch artist locations using the formatted URL
		locations, err := FetchArtistLocations(url)
		if err == nil {
			artists[i].Locations = strings.Join(locations, ", ")
		} else {
			log.Printf("Error fetching location for artist %d: %v", artists[i].ID, err)
		}
	}

	filteredArtists := filterArtists(artists, query)

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

	// Use a buffer to render the template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		ErrorHandler(w, "We encountered an issue while rendering the page. Please try again later.", http.StatusInternalServerError, true, true)
		return
	}

	// Write the rendered template to the response
	_, err = buf.WriteTo(w)
	if err != nil {
		if strings.Contains(err.Error(), "broken pipe") || strings.Contains(err.Error(), "connection reset by peer") {
			log.Println("Client disconnected before response was fully sent")
		} else {
			log.Printf("Error writing response: %v", err)
		}
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

		// locations
		if strings.Contains(strings.ToLower(a.Locations), strings.ToLower(query)) {
			result = append(result, a)
			continue
		}
	}
	return result
}

func GetSearchSuggestionsHandler(w http.ResponseWriter, r *http.Request) {
	initCache()
	query := r.URL.Query().Get("q")
	if query == " " {
		json.NewEncoder(w).Encode([]string{})
		return
	}
	suggestions := []string{}
	artists, _, _, _ := getCachedData()

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

		// First album date
		if strings.Contains(strings.ToLower(artist.FirstAlbum), strings.ToLower(query)) {
			suggestions = append(suggestions, fmt.Sprintf("%s - first album date", artist.FirstAlbum))
		}

		// Creation date
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			suggestions = append(suggestions, fmt.Sprintf("%d - creation date", artist.CreationDate))
		}

		// Locations
		locations := strings.Split(artist.Locations, ", ")
		for _, loc := range locations {
			if strings.Contains(strings.ToLower(loc), strings.ToLower(query)) {
				suggestions = append(suggestions, fmt.Sprintf("%s - location", loc))
			}
		}

	}

	json.NewEncoder(w).Encode(suggestions)
}
// Serve artist details page
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

	artist, location, date, relation, err := api.GetArtistByID(id)
	if err != nil {
		log.Printf("Error retrieving artist by ID %v: %s", id, err)
		ErrorHandler(w, "Ooops!\n We ran into an issue while fetching Artists,\n Please try again later.", http.StatusInternalServerError, false, false)
		return
	}

	data := ArtistDetailData{
		Artist:   *artist,
		Location: *location,
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
	initCache()
	_, locations, _, _ := getCachedData()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		log.Printf("Error encoding locations data to JSON: %v", err)
		ErrorHandler(w, "An error occurred while processing location data. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

// GetDatesHandler handles the /dates route
func GetDatesHandler(w http.ResponseWriter, r *http.Request) {
	initCache()
	_, _, dates, _ := getCachedData()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(dates); err != nil {
		log.Printf("Error encoding dates data to JSON: %v", err)
		ErrorHandler(w, "An error occurred while processing date information. Please try again later.", http.StatusInternalServerError, false, false)
		return
	}
}

// GetRelationsHandler handles the /relations route
func GetRelationsHandler(w http.ResponseWriter, r *http.Request) {
	initCache()
	_, _, _, relations := getCachedData()

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

// Fetch artist locations from the provided URL
func FetchArtistLocations(locationsURL string) ([]string, error) {
	body, err := api.FetchData(locationsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch location data: %v", err)
	}

	var locationData api.Location
	if err := json.Unmarshal(body, &locationData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal location data: %v", err)
	}

	return locationData.Locations, nil
}
