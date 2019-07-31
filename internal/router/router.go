package routers

import (
	"github.com/geofence/internal/controller"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"github.com/geofence/internal/configuration"
)

type WithCORS struct {
	S *mux.Router
}

func (s WithCORS) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}

	// Stop here for a Preflighted OPTIONS request.
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.S.ServeHTTP(res, req)
}

//InitRoutes initialize all routes
func InitRoutes(router WithCORS,
	polyController *controller.PolyController,
	circleController *controller.CircleController,
	appConfig *configuration.Config,
	log log.Logger,
) WithCORS {
	SetGeofencerV1Routes(router.S, *polyController, *circleController)
	return router
}
