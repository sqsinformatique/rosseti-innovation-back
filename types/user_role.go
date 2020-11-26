package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

var (
	ErrBadRole = errors.New("bad role")
)

type Role int

const (
	RestrictedUser Role = iota // default value
	Electrician
	Master
	Engineer
	Admin
)

var stringToRole = map[string]Role{
	"RESTRICTED_USER": RestrictedUser,
	"ELECTRICIAN":     Electrician,
	"MASTER":          Master,
	"ENGINEER":        Engineer,
	"ADMIN":           Admin,
}

func (role Role) String() string {
	for key, item := range stringToRole {
		if item == role {
			return key
		}
	}

	return ""
}

// UnmarshalJSON method is called by json.Unmarshal,
// whenever it is of type Role
func (role *Role) UnmarshalJSON(data []byte) error {
	var roleName string

	if data == nil {
		*role = RestrictedUser
		return nil
	}

	if err := json.Unmarshal(data, &roleName); err != nil {
		return err
	}

	// Check received Role
	if roleName == "" {
		*role = RestrictedUser
	} else {
		r, ok := stringToRole[roleName]
		if !ok {
			return ErrBadRole
		}
		*role = r
	}

	return nil
}

// MarshalJSON method is called by json.Marshal,
// whenever it is of type Role
func (role *Role) MarshalJSON() ([]byte, error) {
	roleName := role.String()

	if roleName == "" {
		return nil, ErrBadRole
	}

	return json.Marshal(roleName)
}

// Value implements the driver Valuer interface.
func (role Role) Value() (driver.Value, error) {
	roleName := role.String()

	if roleName == "" {
		return nil, ErrBadRole
	}

	return roleName, nil
}

// Make the Role struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (role *Role) Scan(value interface{}) error {
	if value == nil {
		*role = RestrictedUser
		return nil
	}

	b, ok := value.(string)

	if !ok {
		return errors.New("type assertion to string failed")
	}

	*role = stringToRole[b]

	return nil
}
