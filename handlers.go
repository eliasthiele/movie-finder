package main

import (
	"net/http"
	"os"
	"log"
	"encoding/json"
	"net/url"
)

func recsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request=%s", r.URL.String())
	token := os.Getenv("TMDB_TOKEN")

	// Receive and handle user params from the http.Request
	urlParams := r.URL.Query()
	log.Printf("urlParams=%s", urlParams)
	params, err := parseParams(urlParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
	// Format a query for TMDB API
	url := "https://api.themoviedb.org/3/discover/movie"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.URL.RawQuery = params.Encode()

	log.Printf("Request url=%s", req.URL.String())

	// Call TMDB API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	
	// Collect and decode the response
	var response DiscoverResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Write the final results
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(resp.StatusCode)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func directorHandler(w http.ResponseWriter, r *http.Request) {
	// Prepare the query
	token := os.Getenv("TMDB_TOKEN")
	nameQuery := r.URL.Query().Get("with_crew")
	log.Printf("Director=%s", nameQuery)
	params := url.Values{}
	params.Set("query", nameQuery)
	url := "https://api.themoviedb.org/3/search/person"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.URL.RawQuery = params.Encode()

	// Get the director info
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Return the director info
	var response DirectorResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(resp.StatusCode)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
