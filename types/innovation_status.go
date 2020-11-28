package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	ErrBadStatus = errors.New("bad status")
)

type Status int

const (
	Unknown Status = iota // default value
	Draft
	Approved
	Expertise
	Revoked
	Recognized
	Experiment
	ExperimentSuccess
	ExperimentFailed
	ReplicationSuccess
	ReplicationFailed
)

var stringToStatus = map[string]Status{
	"UNKNOWN":             Unknown,
	"DRAFT":               Draft,
	"APPROVED":            Approved,
	"EXPERTISE":           Expertise,
	"REVOKED":             Revoked,
	"RECOGNIZED":          Recognized,
	"EXPERIMENT":          Experiment,
	"EXPERIMENT_SUCCESS":  ExperimentSuccess,
	"EXPERIMENT_FAILED":   ExperimentFailed,
	"REPLICATION_SUCCESS": ReplicationSuccess,
	"REPLICATION_FAILED":  ReplicationFailed,
}

func (st Status) String() string {
	for key, item := range stringToStatus {
		if item == st {
			return key
		}
	}

	return ""
}

// UnmarshalJSON method is called by json.Unmarshal,
// whenever it is of type Role
func (st *Status) UnmarshalJSON(data []byte) error {
	var roleName string

	if data == nil {
		*st = Unknown
		return nil
	}

	if err := json.Unmarshal(data, &roleName); err != nil {
		return err
	}

	// Check received Role
	if roleName == "" {
		*st = Unknown
	} else {
		r, ok := stringToStatus[roleName]
		if !ok {
			return ErrBadStatus
		}
		*st = r
	}

	return nil
}

// MarshalJSON method is called by json.Marshal,
// whenever it is of type Role
func (st *Status) MarshalJSON() ([]byte, error) {
	stName := st.String()

	if stName == "" {
		return nil, ErrBadStatus
	}

	return json.Marshal(stName)
}

// Value implements the driver Valuer interface.
func (st Status) Value() (driver.Value, error) {
	stName := st.String()

	if stName == "" {
		return nil, ErrBadStatus
	}

	return stName, nil
}

// Make the Role struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (st *Status) Scan(value interface{}) error {
	if value == nil {
		*st = Unknown
		return nil
	}

	b, ok := value.(string)

	if !ok {
		return errors.New("type assertion to string failed")
	}

	*st = stringToStatus[b]

	return nil
}
