package main

import (
	"log"
	"net/http"
	"os"
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

	addr := ":8080"
	if env, found := os.LookupEnv("ADDR"); found {
		addr = env
	}
	log.Printf("Server on %s running...\n", addr)
	err := http.ListenAndServe(addr, nil)
	log.Fatal(err)
}
