package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Define the directories to scan
const (
	movieDir = "/home/devon/movies"
	tvDir = "/home/devon/tv_shows"
)

// MediaItem represents a media item's name
type MediaItem struct {
	Name string `json:"name"`
}

// scanDirectory scans a directory and returns a list of file names w/o extensions
func scanDirectory(dir string) ([]string, error) {
	var names []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Extract the name w/o extension
			name := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			names = append(names, name)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return names, nil
}

// searchMedia filters a list of names based on a search term (case-insensitive substring match).
func searchMedia(names []string, term string) []MediaItem {
	term = strings.ToLower(term)
	var results []MediaItem
	for _, name := range names {
		if strings.Contains(strings.ToLower(name), term) {
			results = append(results, MediaItem{Name: name})
		}
	}
	return results
}

// movieHandler handles requests to /movie.
func movieHandler(w http.ResponseWriter, r *http.Request) {
	names, err := scanDirectory(movieDir)
	if err != nil {
		http.Error(w, "Failed to scan movies directory", http.StatusInternalServerError)
		return
	}

	term := r.URL.Query().Get("search")
	results := searchMedia(names, term)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// tvHandler handles requests to /tv
func tvHandler(w http.ResponseWriter, r *http.Request) {
	names, err := scanDirectory(tvDir)
	if err != nil {
		http.Error(w, "Failed to scan TV directory", http.StatusInternalServerError)
		return
	}

	term := r.URL.Query().Get("search")
	results := searchMedia(names, term)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func main() {
	// Set up routes.
	http.HandleFunc("/movie", movieHandler)
	http.HandleFunc("/tv", tvHandler)

	// Start the server, log any issues
	fmt.Println("Server starting on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}