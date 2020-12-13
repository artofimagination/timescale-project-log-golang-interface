package timescaledb

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/artofimagination/timescaledb-project-log-go-interface/models"
	"github.com/pkg/errors"
)

var ErrFailedToAdd = errors.New("Failed to add project data")
var ErrFailedToDelete = errors.New("Failed to delete project data")

var AddDataQuery = "INSERT INTO projects_data VALUES (NOW(), ?, ?)"

func (f *TimescaleFunctions) AddData(data []models.Data) error {
	query := AddDataQuery + strings.Repeat(",(NOW(), ?, ?)", len(data)-1) + ")"
	interfaceList := make([]interface{}, len(data))
	for i := range data {
		binary, err := json.Marshal(data[i].Data)
		if err != nil {
			return err
		}
		interfaceList[i] = data[i].ViewerID
		interfaceList[i] = binary
	}

	tx, err := f.Connect()
	if err != nil {
		return err
	}

	result, err := tx.Exec(query, interfaceList...)
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
		return ErrFailedToAdd
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

var GetDataByViewerAndTimeQuery = "SELECT * FROM projects_data WHERE viewer_id = ? AND created_at > ? limit ?"

// GetDataByViewerAndTime returns a chunk of data belonging to the defined viewer and starting from the defined time.
func (f *TimescaleFunctions) GetDataByViewerAndTime(viewerID int, time time.Time, chunkSize int) ([]models.Data, error) {
	tx, err := f.Connect()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(GetDataByViewerAndTimeQuery, viewerID, time, chunkSize)
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
