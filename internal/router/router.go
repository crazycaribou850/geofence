package routers

import (
	"github.com/geofence/internal/controller"
	"github.com/gorilla/mux"
	"log"

	"github.com/geofence/internal/configuration"
)

//InitRoutes initialize all routes
func InitRoutes(router *mux.Router,
	polyController *controller.PolyController,
	circleController *controller.CircleController,
	appConfig *configuration.Config,
	log log.Logger,
) *mux.Router {
	SetGeofencerV1Routes(router, *polyController, *circleController)
	return router
}
