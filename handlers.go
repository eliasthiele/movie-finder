package main

import (
	"net/http"
	"log"
	"encoding/json"
	"net/url"
)

func recsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request=%s", r.URL.String())

	urlParams := r.URL.Query()
	params, err := parseParams(urlParams)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	
	url := "https://api.themoviedb.org/3/discover/movie?"+params.Encode()
	var response DiscoverResponse
	err = callTMDB(url, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if directorID := urlParams.Get("with_crew"); directorID != "" {
		err := filterByDirector(&response, directorID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = addMovieInfo(&response)
	log.Printf("Response: %+v", response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func personHandler(w http.ResponseWriter, r *http.Request) {
	nameQuery := r.URL.Query().Get("person")
	params := url.Values{}
	params.Set("query", nameQuery)
	url := "https://api.themoviedb.org/3/search/person?"+params.Encode()

	var response PersonResponse
	err := callTMDB(url, &response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
