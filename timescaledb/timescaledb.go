package timescaledb

import (
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"

	// Need to register postgres drivers with database/sql
	_ "github.com/lib/pq"
)

var DBConnection = ""
var MigrationDirectory = ""

func BootstrapData() error {
	log.Println("Executing TimeScaleDB migration")

	migrations := &migrate.FileMigrationSource{
		Dir: MigrationDirectory,
	}
	log.Println("Getting migration files")

	db, err := sql.Open("postgres", DBConnection)
	if err != nil {
		return err
	}
	log.Println("DB connection open")

	n := 0
	for retryCount := 5; retryCount > 0; retryCount-- {
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

func ConnectData() (*sql.DB, error) {
	log.Println("Connecting to TimescaleDB")

	db, err := sql.Open("postgres", DBConnection)

	// if there is an error opening the connection, handle it
	if err != nil {
		return nil, err
	}

	return db, nil
}
