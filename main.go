package main

import (
	"net/http"
	"log"
	"fmt"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	http.HandleFunc("/recs", recsHandler)
	http.HandleFunc("/director", directorHandler)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := os.Getenv("PORT")
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
