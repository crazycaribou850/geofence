package main

import (
	"github.com/geofence/router"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}
	router.InitRoutes()
	log.Fatal(http.ListenAndServe(port, nil))
}
