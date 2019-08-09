package routers

import (
	"github.com/gorilla/mux"
	"github.com/geofence/internal/controller"
)

// SetGeofencerV1Routes sets V1 routes
func SetGeofencerV1Routes(router *mux.Router, polyController controller.PolyController, circleController controller.CircleController) {
	polyRouter := router.PathPrefix("/poly").Subrouter()

	polyRouter.Path("/").HandlerFunc(polyController.DetermineMembership()).Methods("POST")
	polyRouter.Path("/all").HandlerFunc(polyController.Ping()).Methods("POST")
	polyRouter.Path("/find/{id}").HandlerFunc(polyController.FindPolyLocationFromID()).Methods("GET")
	polyRouter.Path("/find").HandlerFunc(polyController.FeatureQuery()).Methods("POST")
	polyRouter.Path("/echo").HandlerFunc(polyController.Echo()).Methods("POST", "OPTIONS")
	polyRouter.Path("/closest").HandlerFunc(polyController.FindMostProbableStore()).Methods("POST")
	polyRouter.Path("/intersects").HandlerFunc(polyController.DetermineGeogMembership()).Methods("POST")
	polyRouter.Path("/intersects/{id}").HandlerFunc(polyController.DetermineGeogMembershipFromID()).Methods("POST")

	insertRouter := router.PathPrefix("/insert").Subrouter()
	insertRouter.Path("/poly").HandlerFunc(polyController.InsertPolygon()).Methods("POST")

	circleRouter := router.PathPrefix("/circle").Subrouter()
	circleRouter.Path("/").HandlerFunc(circleController.DetermineMembership()).Methods("POST")
}
