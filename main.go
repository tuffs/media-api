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
	Name			string		`json:"name"`
	Episodes	[]string	`json:"episodes,omitempty"` // Exclusive to TV Shows
}

// nameTidier removes the file extension and 
// then removes the periods underscores and hyphens
// and replaces them with a space for formatted responses
func nameTidier(name string) string {
	// Remove the file extension
	name = strings.TrimSuffix(name, filepath.Ext(name))

	// Replace underscores, hyphens and periods with spaces
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	// Replace multiple periods with a single space
	name = strings.ReplaceAll(name, ".", " ")
	// Trim any extra spaces from the string
	return strings.TrimSpace(name)
}

// scanDirectory scans a directory and returns a list of file names w/o extensions
func scanDirectory(dir string) ([]string, error) {
	var names []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Tidy up the name after removing extension and separators
			tidiedName := nameTidier(info.Name())
			names = append(names, tidiedName)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return names, nil
}

// scanTVDirectory scans the TV directory for show folder and their episodes recursively
func scanTVDirectory(dir, term string) ([]MediaItem, error) {
	var results []MediaItem
	term = strings.ToLower(term)
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Look for directories that might represent TV Shows
		if info.IsDir() && path != dir {
			// Tidy the directory name for matching
			dirName := nameTidier(info.Name())
			if term == "" || strings.Contains(strings.ToLower(dirName), term) {
				// Found a matching show directory, scan for episodes recursively
				var episodes []string
				err := filepath.Walk(path, func(epPath string, epInfo os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !epPath.IsDir() {
						// Collect episode files, tidy their names
						tidiedName := nameTidier(epInfo.Name())
						episodes = append(epidoes, tidiedName)
					}
					return nil
				})
				if err != nil {
					return err
				}
				// Only add the show if it has episodes
				if len(episodes) > 0 {
					results = append(results, MediaItem{
						Name:			dirName,
						Episodes: episodes,
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, nil
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
