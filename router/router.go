package router

import (
	"github.com/geofence/controller"
	"net/http"
	"github.com/gorilla/mux"
)

func InitRoutes(router *mux.Router) *mux.Router {
	http.HandleFunc("/circle/", controller.CircleHandler)
	http.HandleFunc("/poly/", controller.PolyHandler)
}