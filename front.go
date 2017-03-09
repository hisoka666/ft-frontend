package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", index)
	http.Handle("/script/", http.StripPrefix("/script/", http.FileServer(http.Dir("/script"))))
	http.ListenAndServe(":9090", nil)
	log.Println("Listening...")
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
