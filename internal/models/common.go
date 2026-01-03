package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type FileUpload struct {
	URL      string `json:"url"`
	PublicID string `json:"public_id"`
}

// Value implements the driver.Valuer interface for database serialization
func (f FileUpload) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan implements the sql.Scanner interface for database deserialization
func (f *FileUpload) Scan(value interface{}) error {
	if value == nil {
		*f = FileUpload{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		// handle string case (e.g. if driver returns string for JSON/TEXT)
		str, ok := value.(string)
		if !ok {
			return errors.New("type assertion to []byte or string failed")
		}
		bytes = []byte(str)
	}

	return json.Unmarshal(bytes, f)
}
