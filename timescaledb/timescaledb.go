package timescaledb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/artofimagination/timescaledb-project-log-go-interface/models"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"

	// Need to register postgres drivers with database/sql
	_ "github.com/lib/pq"
)

type FunctionsCommon interface {
	BootstrapTables() error
	Connect() (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error

	AddData(data []models.Data) error
	DeleteByViewerID(viewerID int) error
	DeleteByTime(intervalString string) error
	GetDataByViewerAndTime(viewerID int, time time.Time, chunkSize int) ([]models.Data, error)
}

type TimescaleFunctions struct {
	DBConnection       string
	MigrationDirectory string
}

func (*TimescaleFunctions) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (*TimescaleFunctions) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (c *TimescaleFunctions) BootstrapTables() error {
	log.Println("Executing TimeScaleDB migration")
	migrations := &migrate.FileMigrationSource{
		Dir: c.MigrationDirectory,
	}

	db, err := sql.Open("postgres", c.DBConnection)
	if err != nil {
		return err
	}

	n := 0
	for retryCount := 15; retryCount > 0; retryCount-- {
		n, err = migrate.Exec(db, "postgres", migrations, migrate.Up)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		log.Printf("Failed to execute migration %s. Retrying...\n", err.Error())
	}

	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Migration failed after multiple retries.")
	}
	log.Printf("Applied %d migrations!\n", n)
	return nil
}

func (*TimescaleFunctions) RollbackWithErrorStack(tx *sql.Tx, errorStack error) error {
	if err := tx.Rollback(); err != nil {
		errorString := fmt.Sprintf("%s\n%s\n", errorStack.Error(), err.Error())
		return errors.Wrap(errors.WithStack(errors.New(errorString)), "Failed to rollback changes")
	}
	return errorStack
}

func (c *TimescaleFunctions) Connect() (*sql.Tx, error) {
	db, err := sql.Open("postgres", c.DBConnection)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}
