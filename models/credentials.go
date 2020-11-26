package models

import (
	// local
	"errors"

	"github.com/sqsinformatique/rosseti-innovation-back/types"
)

var (
	ErrEmptyCredentials = errors.New("empty credentials")
)

// Credentials is a struct for check user
type Credentials struct {
	Password string `json:"user_password"`
	Email    string `json:"user_email"`
	Phone    string `json:"user_phone"`
}

func (c *Credentials) String() string {
	return "Mail: " + c.Email + ", Login: " + c.Phone
}

func (c *Credentials) Validate() error {
	if c.Password == "" || (c.Phone == "" && c.Email == "") {
		return ErrEmptyCredentials
	}

	return nil
}

// NewCredentials is a stuct for create user
type NewCredentials struct {
	Password string     `json:"user_password"`
	Email    string     `json:"user_email"`
	Phone    string     `json:"user_phone"`
	Role     types.Role `json:"user_role"`
}

func (c *NewCredentials) String() string {
	return "Email: " + c.Email + ", Phone: " + c.Phone[0:4] + "****" + c.Phone[len(c.Phone)-2:] + ", Role: " + c.Role.String()
}

func (c *NewCredentials) Validate() error {
	if c.Password == "" || c.Phone == "" || c.Email == "" {
		return ErrEmptyCredentials
	}

	return nil
}

// UpdateCredentials is a structs for update user's Credentials
type UpdateCredentials struct {
	Password    string `json:"user_password"`
	OldPassword string `json:"user_old_password"`
}

func (c *UpdateCredentials) Validate() error {
	if c.Password == "" {
		return ErrEmptyCredentials
	}
	return nil
}
