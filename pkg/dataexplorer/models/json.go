package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Json json.RawMessage

func (j *Json) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = Json(result)
	return err
}

func (j Json) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}
