package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullString is a wrapper around sql.NullString
type NullTime struct {
	sql.NullTime
}

// MarshalJSON method is called by json.Marshal,
// whenever it is of type NullString
func (x *NullTime) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(x.Time)
}

// UnmarshalJSON method is called by json.Unmarshal,
// whenever it is of type NullTime
func (x *NullTime) UnmarshalJSON(data []byte) error {
	var v time.Time

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v.Year() == 1 {
		return nil
	}

	x.Time = v
	x.Valid = true

	return nil
}

// Scan implements the Scanner interface.
func (x *NullTime) Scan(value interface{}) error {
	return x.NullTime.Scan(value)
}

// Value implements the driver Valuer interface.
func (x NullTime) Value() (driver.Value, error) {
	return x.NullTime.Value()
}
