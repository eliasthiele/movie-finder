package main

import (
	"net/http"
	"log"
	"fmt"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/recs", recsHandler)
	http.HandleFunc("/director", directorHandler)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Listening on port 9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
