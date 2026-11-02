package main

import (
	"net/url"
	"log"
	"strings"
	"fmt"
	"github.com/biter777/countries"
)

func parseParams(p url.Values) (url.Values, error) {
	params := url.Values{}
	if genres := p.Get("with_genres"); genres != "" {
		log.Println(genres)
		genreIds, err := replaceGenres(genres)
		if err != nil {
			return nil, err
		}
		params.Add("with_genres", genreIds)
	}
	if origin := p.Get("with_origin_country"); origin != "" {
		originCode, err := replaceCountry(origin)
		if err != nil {
			return nil, err
		}
		params.Add("with_origin_country", originCode)
	}
	return params, nil
}

func replaceGenres(s string) (string, error) {
	log.Printf("Genre string received=%s", s)
	genres := strings.Split(s, ",")
	ids := make([]string, 0, len(genres))
	for _, genre := range genres {
		genre = strings.ToLower(strings.TrimSpace(genre))
		id, ok := genreToId[genre]
		if !ok {
			return "", fmt.Errorf("unknown genre: %s", genre)
		}
		ids = append(ids, id)
	}
	return strings.Join(ids, ","), nil
}

func replaceCountry(s string) (string, error) {
	c := countries.ByName(s)
	if c == countries.Unknown {
		return "", fmt.Errorf("unknown country: %s", s)
	}
	return c.Alpha2(), nil
}
