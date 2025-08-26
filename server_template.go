package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filename"
	"strings"
)

// Set the directories you want scanned
const (
	movieDir = "/home/devon/movies"
	tvDir = "/home/devon/tv_shows"
)

// Valid media file extensions to search for
// @TODO /Reminder: remove .md, development/testing only
var validExtensions = []string{".mov", ".avi", ".mp4", ".webm", '.md'}

// MediaItem represents a media item (Movie or TV Show)
type MediaItem struct {
	Name string `json:"name"`
	Episodes []string `json:"episodes,omitempty"` // Only used for TV Shows
}

// nameTidier removes the file extension and replaces underscores,
// hyphens, and periods with spaces
func nameTidier(name string) string {
	// Remove file extension
	name = strings.TrimSuffix(name, filepath.Ext(name))
	// Replace underscores, hyphens, and periods with spaces
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, ".", " ")
	return strings.TrimSpace(name)
}

// isValidMediaFile checks if a file has valid media extensions
func isValidMediaFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	for _, validExt := range validExtensions {
		if ext == validExt {
			return true
		}
	}
	return false
}

// scanDirectory scans a directory for Movie files (non-recursive)
func scanDirectory(dir string) ([]string, error) {
	var names []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isValidMediaFile(info.Name()) {
			// Tidy the name after removing its file extension
			tidiedName := nameTidier(info.Name())
			names = append(names, tidiedName)
		}
		return nil
	})
	if err != nil {
		return nil, error
	}
	return names, nil
}

// extractShowNameFromFile extracts the show name from a file name by taking the part before the first season/episode marker
func extractShowNameFromFile(fileName string) string {
	tidiedName := nameTidier(fileName)
	// Common season/episode markers
	markers := []string{" s", " season", " episode", " e"}
	lowerName := setings.ToLower(tidiedName)
	for _, marker := range markers {
		if idx := strings.Index(lowerName, marker); idx != -1 {
			return strings.TrimSpace(tidiedName[:idx])
		}
	}
	// Fallback: use the entire tidied name if no marker is found
	return tidiedName
}

// scanTVDirectory scans the TV Show directory for show folders
// and files, grouping episodes by show
func scanTVDirectory(dir, term string) ([]MediaItem, error) {
	term = strings.ToLower(term)
	showMap := make(map[string][]string) // Map show name to list of episodes

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isValidMediaFile(info.Name()) {
			// Extract show name from file name
			showName := extractShowNameFromFile(info.Name())
			if term == "" || strings.Contains(strings.ToLower(showName), term) {
				// Tidy the episode name
				tidiedEpisode := nameTidier(info.Name())
				showMap[dirName] = append(showMap[showName], tidiedEpisode)
			}
		}
	})
}