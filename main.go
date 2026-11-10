package main

import (
	"net/http"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/recs", recsHandler)
	http.HandleFunc("/person", personHandler)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}
	address := "0.0.0.0:"+port
	log.Printf("Listening on port %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
