package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Umleitung von / auf /app/
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/app/")
	})

	// Statische Dateien servieren
	r.StaticFS("/app/", http.Dir("./static/app/"))

	// Einstellungen für den Server-Adresse über Umgebungsvariable
	addr := ":8080"
	if env, found := os.LookupEnv("ADDR"); found {
		addr = env
	}

	// Server starten
	log.Printf("Server on %s running...\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
