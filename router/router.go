package router

import (
	"github.com/geofence/controller"
	"net/http"
)

func InitRoutes() {
	http.HandleFunc("/circle/", controller.CircleHandler)
	http.HandleFunc("/poly/", controller.PolyHandler)
}