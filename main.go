package main

import (
	"github.com/geofence/router"
	"log"
	"net/http"

)

func main() {
	router.InitRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
