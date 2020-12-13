package dbcontrollers

import (
	"fmt"
	"os"

	"github.com/artofimagination/timescaledb-project-log-go-interface/timescaledb"
	"github.com/pkg/errors"
)

type DBControllerCommon interface {
}

type TimescaleController struct {
	DBFunctions timescaledb.FunctionsCommon
}

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

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		pass,
		address,
		port,
		dbName)
	migrationDirectory := os.Getenv("TIMESCALE_DB_MIGRATION_DIR")
	if migrationDirectory == "" {
		return nil, errors.New("TIMESCALE DB migration folder not defined")
	}

	controller := &TimescaleController{
		DBFunctions: &timescaledb.TimescaleFunctions{
			MigrationDirectory: migrationDirectory,
			DBConnection:       connectionString,
		},
	}

	if err := controller.DBFunctions.(*timescaledb.TimescaleFunctions).BootstrapTables(); err != nil {
		return nil, fmt.Errorf("Data bootstrap failed. %s", errors.WithStack(err))
	}

	return controller, nil
}
