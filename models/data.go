package models

import (
	"time"
)

type Data struct {
	CreatedAt time.Time `validation:"required"`
	ViewerID  int       `validation:"required"`
	Data      DataMap   `validation:"required"`
}

type DataMap map[string]interface{}
