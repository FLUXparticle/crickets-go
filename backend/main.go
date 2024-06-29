package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("Server on :8080 running...")
	err := http.ListenAndServe("localhost:8080", nil)
	log.Fatal(err)
}
