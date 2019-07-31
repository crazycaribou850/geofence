package application

import (
	"github.com/geofence/internal/configuration"
	"github.com/geofence/internal/controller"
	"github.com/geofence/internal/db"
	r "github.com/geofence/internal/router"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
)

type App struct {
	Port string
	DB *sqlx.DB
	Router r.WithCORS
}

func NewApplication(appConfig *configuration.Config) (*App, error) {
	logger := log.Logger{}
	logger.SetOutput(os.Stdout)
	db, err := db.NewDB(appConfig.DBURL, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating postgres client")
	}

	polyController := controller.NewPolyController(validator.New(), logger, db)
	circleController := controller.NewCircleController(validator.New(), logger)
	router := r.WithCORS{mux.NewRouter()}
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
