package routers

import (
	"github.com/gorilla/mux"
	"github.com/geofence/internal/controller"
)

// SetGeofencerV1Routes sets V1 routes
func SetGeofencerV1Routes(router *mux.Router, polyController controller.PolyController, circleController controller.CircleController) {
	polyRouter := router.PathPrefix("/poly").Subrouter()
	polyRouter.Path("/").HandlerFunc(polyController.DetermineMembership()).Methods("POST")
	circleRouter := router.PathPrefix("/circle").Subrouter()
	circleRouter.Path("/").HandlerFunc(circleController.DetermineMembership()).Methods("POST")
}
