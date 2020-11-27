package models

import "github.com/sqsinformatique/rosseti-innovation-back/types"

type Innovation struct {
	ID          int            `json:"id" db:"id"`
	AuthorID    int            `json:"author_id" db:"author_id"`
	Title       string         `json:"title" db:"title"`
	Tags        string         `json:"tags" db:"tags"`
	Description string         `json:"descriptions" db:"descriptions"`
	State       string         `json:"state" db:"state"`
	Meta        types.NullMeta `json:"meta" db:"meta"`
	Timestamp
}

func (u *Innovation) SQLParamsRequest() []string {
	return []string{
		"author_id",
		"title",
		"tags",
		"descriptions",
		"state",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

type InnovationCoAuthors struct {
	ID       int `json:"id" db:"id"`
	AuthorID int `json:"author_id" db:"author_id"`
	Timestamp
}

func (u *InnovationCoAuthors) SQLParamsRequest() []string {
	return []string{
		"id",
		"author_id",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

type InnovationExperts struct {
	ID       int `json:"id" db:"id"`
	ExpertID int `json:"expert_id" db:"expert_id"`
	Timestamp
}

func (u *InnovationExperts) SQLParamsRequest() []string {
	return []string{
		"id",
		"expert_id",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

type InnovationFiles struct {
	ID       int    `json:"id" db:"id"`
	FileID   string `json:"file_id" db:"file_id"`
	FileName string `json:"file_name" db:"file_name"`
	Timestamp
}

func (u *InnovationFiles) SQLParamsRequest() []string {
	return []string{
		"id",
		"file_id",
		"file_name",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}
