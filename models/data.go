package models

import (
	"time"

	"github.com/google/uuid"
)

type Data struct {
	CreatedAt time.Time `validation:"required"`
	ViewerID  uuid.UUID `validation:"required"`
	Data      DataMap   `validation:"required"`
}

type DataMap map[string]interface{}
