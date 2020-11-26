package models

import (
	"github.com/sqsinformatique/rosseti-innovation-back/types"
)

type Profile struct {
	ID               int            `json:"id" db:"id"`
	FirstName        string         `json:"user_first_name" db:"user_first_name"`
	MiddleName       string         `json:"user_middle_name" db:"user_middle_name"`
	LastName         string         `json:"user_last_name" db:"user_last_name"`
	Position         string         `json:"user_position" db:"user_position"`
	Company          string         `json:"user_company" db:"user_company"`
	UserElectroGroup string         `json:"user_electro_group" db:"user_electro_group"`
	PrivateKey       string         `json:"-" db:"user_private_key"`
	PublicKey        string         `json:"user_public_key" db:"user_public_key"`
	Meta             types.NullMeta `json:"meta" db:"meta"`
	Timestamp
}

func (u *Profile) SQLParamsRequest() []string {
	return []string{
		"id",
		"user_first_name",
		"user_middle_name",
		"user_last_name",
		"user_position",
		"user_company",
		"user_electro_group",
		"user_private_key",
		"user_public_key",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}
