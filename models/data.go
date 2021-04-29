package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Data struct {
	CreatedAt time.Time `json:"created_at" validation:"required"`
	ViewerID  uuid.UUID `json:"viewer_id" validation:"required"`
	Data      DataMap   `json:"data" validation:"required"`
}

type DataMap map[string]interface{}

func (a DataMap) Value() (driver.Value, error) {
	return json.Marshal(a)
}
