package main

import (
	"net/http"
	"encoding/json"
	"net/url"
	//"log"
	"strings"
	"fmt"
	"github.com/biter777/countries"
	"os"
)

func parseParams(p url.Values) (url.Values, error) {
	params := url.Values{}
	if genres := p.Get("with_genres"); genres != "" {
		genreIDs, err := replaceGenres(genres)
		if err != nil { return nil, err }
		params.Add("with_genres", genreIDs)
	}
	if origin := p.Get("with_origin_country"); origin != "" {
		originCode, err := replaceCountry(origin)
		if err != nil { return nil, err }
		params.Add("with_origin_country", originCode)
	}
	if directorID := p.Get("with_crew"); directorID != "" {
		params.Add("with_crew", directorID)
	}
	if actorID := p.Get("with_cast"); actorID != "" {
		params.Add("with_cast", actorID)
	}
	if primaryReleaseYear := p.Get("primary_release_year"); primaryReleaseYear != "" {
		params.Add("primary_release_year", primaryReleaseYear)
	}
	if maxRuntime := p.Get("maxRuntime"); maxRuntime != "" {
		params.Add("with_runtime.lte", maxRuntime)
	}
	if minRuntime := p.Get("minRuntime"); minRuntime != "" {
		params.Add("with_runtime.gte", minRuntime)
	}
	return params, nil
}

func replaceGenres(s string) (string, error) {
	genres := strings.Split(s, ",")
	ids := make([]string, 0, len(genres))
	for _, genre := range genres {
		genre = strings.ToLower(strings.TrimSpace(genre))
		id, ok := genreToId[genre]
		if !ok { return "", fmt.Errorf("unknown genre: %s", genre) }
		ids = append(ids, id)
	}
	return strings.Join(ids, ","), nil
}

func replaceCountry(s string) (string, error) {
	c := countries.ByName(s)
	if c == countries.Unknown { return "", fmt.Errorf("unknown country: %s", s) }
	return c.Alpha2(), nil
}

func filterByDirector(dr *DiscoverResponse, directorID string) error {
	url := fmt.Sprintf("https://api.themoviedb.org/3/person/%s/movie_credits", directorID)

	var response CreditsResponse
	err := callTMDB(url, &response)
	if err != nil { return err }

	directedFilmIDs := make(map[int]struct{})
	for _, credit := range response.Credits {
		if credit.Job == "Director" {
			directedFilmIDs[credit.ID] = struct{}{}
		}
	}

	filteredMovies := make([]Movie, 0)
	for _, movie := range dr.Results {
		if _, ok := directedFilmIDs[movie.ID]; ok {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	dr.Results = filteredMovies

	return nil
}

// out should be a pointer to a struct
func callTMDB(url string, out any) error {
	token := os.Getenv("TMDB_TOKEN")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil { return err }
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(out)
	if err != nil { return err }

	return nil
}
