package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchData(t *testing.T) {
	mockResponse := `{"id": 1, "name": "Queen"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	data, err := FetchData(server.URL)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := mockResponse
	if string(data) != expected {
		t.Errorf("Expected %v, got %v", expected, string(data))
	}
}

func TestGetArtists(t *testing.T) {
	mockArtists := []Artist{
		{
			ID:           1,
			Name:         "Queen",
			Image:        "https://groupietrackers.herokuapp.com/api/images/queen.jpeg",
			CreationDate: 1970,
			FirstAlbum:   "14-12-1973",
			Members:      []string{"Freddie Mercury", "Brian May", "John Daecon", "Roger Meddows-Taylor", "Mike Grose", "Barry Mitchell", "Doug Fogie"},
			Locations:    "https://groupietrackers.herokuapp.com/api/locations/1",
			ConcertDates: "https://groupietrackers.herokuapp.com/api/dates/1",
			Relations:    "https://groupietrackers.herokuapp.com/api/relation/1",
		},
	}
	mockResponse, _ := json.Marshal(mockArtists)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponse)
	}))
	defer server.Close()

	ArtistsURL = server.URL

	artists, err := GetArtists()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if artists[0].Name != "Queen" {
		t.Errorf("Expected artist names to be 'Queen' , got %v", artists[0].Name)
	}
}

func TestGetLocations(t *testing.T) {
	mockLocations := struct {
		Index []Location `json:"index"`
	}{
		Index: []Location{
			{
				ID:        1,
				Locations: []string{"north_carolina-usa", "georgia-usa", "los_angeles-usa", "saitama-japan", "osaka-japan", "nagoya-japan", "penrose-new_zealand", "dunedin-new_zealand"},
				Dates:     "https://groupietrackers.herokuapp.com/api/dates/1",
			},
		},
	}
	mockResponse, _ := json.Marshal(mockLocations)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponse)
	}))
	defer server.Close()

	LocationsURL = server.URL

	locations, err := GetLocations()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if locations[0].Locations[0] != "north_carolina-usa" {
		t.Errorf("Expected locations to be 'north_carolina-usa', got %v", locations[0].Locations[0])
	}
}

func TestGetDates(t *testing.T) {
	mockDates := struct {
		Index []Date `json:"index"`
	}{
		Index: []Date{
			{
				ID:    1,
				Dates: []string{"*23-08-2019", "*22-08-2019", "*20-08-2019", "*26-01-2020", "*28-01-2020", "*30-01-2019", "*07-02-2020", "*10-02-2020"},
			},
		},
	}
	mockResponse, _ := json.Marshal(mockDates)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponse)
	}))
	defer server.Close()

	DatesURL = server.URL

	dates, err := GetDates()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if dates[0].Dates[0] != "*23-08-2019" {
		t.Errorf("Expected dates to be '*23-08-2019', got %v", dates[0].Dates[0])
	}
}

func TestGetRelations(t *testing.T) {
	mockRelations := struct {
		Index []Relation `json:"index"`
	}{
		Index: []Relation{
			{
				ID: 1,
				DatesLocations: map[string][]string{
					"north_carolina-usa":  {"23-08-2019"},
					"georgia-usa":         {"22-08-2019"},
					"los_angeles-usa":     {"20-08-2019"},
					"saitama-japan":       {"26-01-2020"},
					"osaka-japan":         {"28-01-2020"},
					"nagoya-japan":        {"30-01-2019"},
					"penrose-new_zealand": {"07-02-2020"},
					"dunedin-new_zealand": {"10-02-2020"},
				},
			},
		},
	}
	mockResponse, _ := json.Marshal(mockRelations)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(mockResponse)
	}))
	defer server.Close()

	RelationURL = server.URL

	relations, err := GetRelations()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if relations[0].DatesLocations["north_carolina-usa"][0] != "23-08-2019" {
		t.Errorf("Expected relations to have specific dates, got %v", relations[0].DatesLocations["north_carolina-usa"][0])
	}
}

func TestGetArtistsEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	ArtistsURL = server.URL

	artists, err := GetArtists()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(artists) != 0 {
		t.Errorf("Expected empty slice, got %v", artists)
	}
}

func TestGetLocationsError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	LocationsURL = server.URL

	_, err := GetLocations()
	if err == nil {
		t.Fatalf("Expected an error, got nil")
	}
}

func TestGetDatesInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json}"))
	}))
	defer server.Close()

	DatesURL = server.URL

	_, err := GetDates()
	if err == nil {
		t.Fatalf("Expected an error for invalid JSON, got nil")
	}
}

func TestGetRelationsEmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"index":[]}`))
	}))
	defer server.Close()

	RelationURL = server.URL

	relations, err := GetRelations()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(relations) != 0 {
		t.Errorf("Expected empty slice, got %v", relations)
	}
}

func TestGetArtistByID(t *testing.T) {
	setupMockServers := func() (*httptest.Server, *httptest.Server, *httptest.Server, *httptest.Server) {
		artistServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode([]Artist{{ID: 1, Name: "Test Artist"}})
		}))
		locationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(struct{ Index []Location }{[]Location{{ID: 1, Locations: []string{"Test Location"}}}})
		}))
		dateServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(struct{ Index []Date }{[]Date{{ID: 1, Dates: []string{"2023-01-01"}}}})
		}))
		relationServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			json.NewEncoder(w).Encode(struct{ Index []Relation }{[]Relation{{ID: 1, DatesLocations: map[string][]string{"Test Location": {"2023-01-01"}}}}})
		}))
		return artistServer, locationServer, dateServer, relationServer
	}

	artistServer, locationServer, dateServer, relationServer := setupMockServers()
	defer artistServer.Close()
	defer locationServer.Close()
	defer dateServer.Close()
	defer relationServer.Close()

	ArtistsURL = artistServer.URL
	LocationsURL = locationServer.URL
	DatesURL = dateServer.URL
	RelationURL = relationServer.URL

	artist, location, date, relation, err := GetArtistByID(1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if artist == nil || location == nil || date == nil || relation == nil {
		t.Fatalf("Expected all fields to be non-nil")
	}
	if artist.Name != "Test Artist" {
		t.Errorf("Expected artist name to be 'Test Artist', got %v", artist.Name)
	}
}

func TestGetArtistByIDNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]Artist{})
	}))
	defer server.Close()

	ArtistsURL = server.URL
	LocationsURL = server.URL
	DatesURL = server.URL
	RelationURL = server.URL

	_, _, _, _, err := GetArtistByID(999)
	if err == nil {
		t.Fatalf("Expected an error for non-existent artist, got nil")
	}
}