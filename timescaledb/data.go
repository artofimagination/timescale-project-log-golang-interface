package timescaledb

import (
	"encoding/json"
	"time"

	"github.com/artofimagination/timescaledb-project-log-go-interface/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

var ErrFailedToAdd = errors.New("Failed to add project data")
var ErrFailedToDelete = errors.New("Failed to delete project data")

func (f *TimescaleFunctions) AddData(data []models.Data) error {
	tx, err := f.Connect()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("projects_data", "created_at", "viewer_id", "data"))
	if err != nil {
		return err
	}

	for _, value := range data {
		_, err = stmt.Exec(value.CreatedAt, value.ViewerID, value.Data)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return tx.Commit()
}

var DeleteByViewerIDQuery = "DELETE FROM projects_data WHERE viewer_id=?"

func (f *TimescaleFunctions) DeleteByViewerID(viewerID int) error {
	tx, err := f.Connect()
	if err != nil {
		return err
	}

	result, err := tx.Exec(DeleteByViewerIDQuery, viewerID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return f.RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrFailedToDelete
	}

	return tx.Commit()
}

var DeleteByTimeQuery = "SELECT drop_chunks(INTERVAL '?', 'projects_data');"

// DeleteByTime is handling cleanup delete of old data after it has been backed up.
// It will delete all data in the table older, than the INTERVAL.
func (f *TimescaleFunctions) DeleteByTime(intervalString string) error {
	tx, err := f.Connect()
	if err != nil {
		return err
	}

	result, err := tx.Exec(DeleteByTimeQuery, intervalString)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return f.RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		if errRb := tx.Rollback(); errRb != nil {
			return err
		}
		return ErrFailedToDelete
	}

	return tx.Commit()
}

var GetDataByViewerAndTimeQuery = "SELECT * FROM projects_data WHERE viewer_id = $1 AND created_at > $2"

// GetDataByViewerAndTime returns a chunk of data belonging to the defined viewer and starting from the defined time.
func (f *TimescaleFunctions) GetDataByViewerAndTime(viewerID *uuid.UUID, time time.Time) ([]models.Data, error) {
	tx, err := f.Connect()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(GetDataByViewerAndTimeQuery, viewerID, time)
	if err != nil {
		return nil, err
	}

	dataList := make([]models.Data, 0)
	defer rows.Close()
	for rows.Next() {
		data := models.Data{}
		dataMap := []byte{}
		err = rows.Scan(&data.CreatedAt, &data.ViewerID, &dataMap)
		if err != nil {
			return nil, f.RollbackWithErrorStack(tx, err)
		}
		if err := json.Unmarshal(dataMap, &data.Data); err != nil {
			return nil, f.RollbackWithErrorStack(tx, err)
		}
		dataList = append(dataList, data)
	}

	err = rows.Err()
	if err != nil {
		return nil, f.RollbackWithErrorStack(tx, err)
	}

	return dataList, tx.Commit()
}
