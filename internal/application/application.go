package application

import (
	"database/sql"
	"github.com/geofence/internal/configuration"
	"github.com/geofence/internal/controller"
	"github.com/geofence/internal/db"
	"github.com/gorilla/mux"
	r "github.com/geofence/internal/router"
	"github.com/pkg/errors"
	"net/http"
	"log"
	"gopkg.in/go-playground/validator.v9"
)

type App struct {
	Port string
	DB *sql.DB
	Router *mux.Router
}

func NewApplication(appConfig *configuration.Config) (*App, error) {
	db, err := db.NewDB(appConfig.DBURL)
	if err != nil {
		return nil, errors.Wrap(err, "error creating postgres client")
	}
	logger := log.Logger{}
	polyController := controller.NewPolyController(validator.New(), logger)
	circleController := controller.NewCircleController(validator.New(), logger)
	router := mux.NewRouter()
	router = r.InitRoutes(router, polyController, circleController, appConfig, logger)
	return &App{
		Port: appConfig.Port,
		DB: db,
		Router: router,
	}, nil
}


func (a *App) Run() {
	a.Start()
}

func (a *App) Start() {
	http.ListenAndServe(a.Port, a.Router)
}
