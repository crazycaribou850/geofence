package main

import (
	"github.com/geofence/internal/application"
	"github.com/geofence/internal/configuration"
	"github.com/pkg/errors"
	"log"
)

func main() {
	appConfig := configuration.Load()

	app, err := application.NewApplication(appConfig)
	if wErr := errors.Wrapf(err, "failed setting up application"); wErr != nil {
		log.Panic(wErr)
	}
	defer app.DB.Close()

	app.Run()
}