package models

import (
	"errors"

	"github.com/sqsinformatique/rosseti-innovation-back/types"
)

type User struct {
	ID    int            `json:"id" db:"id"`
	Hash  string         `json:"-" db:"user_hash"`
	Role  types.Role     `json:"user_role" db:"user_role"`
	Email string         `json:"user_email" db:"user_email"`
	Phone string         `json:"user_phone" db:"user_phone"`
	Meta  types.NullMeta `json:"meta" db:"meta"`
	Timestamp
}

func (u *User) SQLParamsRequest() []string {
	return []string{
		"user_hash",
		"user_email",
		"user_phone",
		"user_role",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

var (
	ErrEmptyEmail = errors.New("empty email")
	ErrEmptyPhone = errors.New("empty phone")
)

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmptyEmail
	}

	if u.Phone == "" {
		return ErrEmptyPhone
	}

	return nil
}
