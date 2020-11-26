package models

import (
	"github.com/sqsinformatique/rosseti-innovation-back/types"
)

type Session struct {
	ID     string         `json:"id" db:"id"`
	UserID int            `json:"user_id" db:"user_id"`
	Meta   types.NullMeta `json:"meta" db:"meta"`
	Timestamp
}

func (s *Session) SQLParamsRequest() []string {
	return []string{
		"id",
		"user_id",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}
