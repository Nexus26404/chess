package main

import (
	"log"

	"chess/server/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Load routes from the router layer
	router.RegisterRoutes(r)

	// Start HTTPS server
	if err := r.RunTLS(":8443", "certs/server.crt", "certs/server.key"); err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
}
