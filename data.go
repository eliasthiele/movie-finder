package main 

var genreToId = map[string]string{
    "action":          "28",
    "adventure":       "12",
    "animation":       "16",
    "comedy":          "35",
    "crime":           "80",
    "documentary":     "99",
    "drama":           "18",
    "family":          "10751",
    "fantasy":         "14",
    "history":         "36",
    "horror":          "27",
    "music":           "10402",
    "mystery":         "9648",
    "romance":         "10749",
    "science fiction": "878",
    "tv movie":        "10770",
    "thriller":        "53",
    "war":             "10752",
    "western":         "37",
}

type DiscoverResponse struct {
    Page         int     `json:"page"`
    Results      []Movie `json:"results"`
    TotalPages   int     `json:"total_pages"`
    TotalResults int     `json:"total_results"`
}

type Movie struct {
    ID          int     `json:"id"`
    Title       string  `json:"title"`
    ReleaseDate string  `json:"release_date"`  
    PosterPath  string  `json:"poster_path"`
    VoteAverage float64 `json:"vote_average"`
    Synopsis    string  `json:"overview"`
}

type DirectorResponse struct {
    Results []Director `json:"results"`
}

type Director struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
