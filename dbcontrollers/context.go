package dbcontrollers

import (
	"errors"
	"fmt"
	"os"

	"github.com/artofimagination/timescaledb-project-log-go-interface/timescaledb"
)

type DBControllerCommon interface {
}

type TimescaleController struct{}

func NewDBController() (*TimescaleController, error) {
	address := os.Getenv("TIMESCALE_DB_ADDRESS")
	if address == "" {
		return nil, errors.New("TIMESCALE DB address not defined")
	}
	port := os.Getenv("TIMESCALE_DB_PORT")
	if address == "" {
		return nil, errors.New("TIMESCALE DB port not defined")
	}
	username := os.Getenv("TIMESCALE_DB_USER")
	if address == "" {
		return nil, errors.New("TIMESCALE DB username not defined")
	}
	pass := os.Getenv("TIMESCALE_DB_PASSWORD")
	if address == "" {
		return nil, errors.New("TIMESCALE DB password not defined")
	}
	dbName := os.Getenv("TIMESCALE_DB_NAME")
	if address == "" {
		return nil, errors.New("TIMESCALE DB name not defined")
	}

	timescaledb.MigrationDirectory = os.Getenv("TIMESCALE_DB_PASSWORD")
	if timescaledb.MigrationDirectory == "" {
		return nil, errors.New("TIMESCALE DB migration folder not defined")
	}
	timescaledb.DBConnection = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		pass,
		address,
		port,
		dbName)

	controller := &TimescaleController{}
	return controller, nil
}
