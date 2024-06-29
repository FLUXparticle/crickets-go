package main

import (
	"log"
	"net/http"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", loggingMiddleware(fs))

	log.Println("Server on :8080 running...")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
