package main

import (
	"log"
	"net/http"
	"strings"

	"learn.zone01kisumu.ke/git/johnodhiambo0/groupie-tracker/controllers"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/artists", controllers.ServeArtists)
	http.HandleFunc("/artist/", controllers.ServeArtistDetails)
	http.HandleFunc("/about", controllers.AboutHandler)

	// Catch-all for undefined routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is for the root path or favicon.ico
		if r.URL.Path == "/" || r.URL.Path == "/favicon.ico" {
			controllers.ServeArtists(w, r)
			return
		}

		// Check if the request is for an artist details page
		if strings.HasPrefix(r.URL.Path, "/artist/") {
			controllers.ServeArtistDetails(w, r)
			return
		}

		// If it's not a known route, use ErrorHandler for 404
		controllers.ErrorHandler(w, "Page Not Found", http.StatusNotFound, true, true)
	})

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
