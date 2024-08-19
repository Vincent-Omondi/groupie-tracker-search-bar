package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
    // Change to the directory where the templates are located
    if err := os.Chdir(".."); err != nil {
        panic(err)
    }
    os.Exit(m.Run())
}

func TestServeArtists(t *testing.T) {
	req, err := http.NewRequest("GET", "/artists", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeArtists)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Queen") { 
		t.Errorf("Expected response body to contain 'Queen'")
	}
}

func TestServeArtistDetails(t *testing.T) {
	req, err := http.NewRequest("GET", "/artist/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ServeArtistDetails)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Freddie Mercury") { // Change to match expected template output
		t.Errorf("Expected response body to contain 'Freddie Mercury'")
	}
}

func TestGetArtistsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/artists?query=Queen", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetArtistsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Queen") {
		t.Errorf("Expected response body to contain 'Queen'")
	}
}

func TestGetLocationsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/locations", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetLocationsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "north_carolina-usa") {
		t.Errorf("Expected response body to contain 'north_carolina-usa'")
	}
}

func TestGetDatesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/dates", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetDatesHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "*23-08-2019") {
		t.Errorf("Expected response body to contain '*23-08-2019'")
	}
}

func TestGetRelationsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/relations", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetRelationsHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "north_carolina-usa") {
		t.Errorf("Expected response body to contain 'north_carolina-usa'")
	}
}

func TestGetArtistByIDHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/artists/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetArtistByIDHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Queen") {
		t.Errorf("Expected response body to contain 'Queen'")
	}
}

func TestAboutHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/about", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AboutHandler)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK; got %v", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "About") { // Adjust based on the actual content of your About page
		t.Errorf("Expected response body to contain 'About'")
	}
}
