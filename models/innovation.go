package models

import "github.com/sqsinformatique/rosseti-innovation-back/types"

type Innovation struct {
	ID          int            `json:"id" db:"id"`
	AuthorID    int            `json:"author_id" db:"author_id"`
	Title       string         `json:"title" db:"title"`
	Tags        string         `json:"tags" db:"tags"`
	Problem     string         `json:"problem" db:"problem"`
	Description string         `json:"descriptions" db:"descriptions"`
	State       types.Status   `json:"state" db:"state"`
	Meta        types.NullMeta `json:"meta" db:"meta"`
	Timestamp
}

func (u *Innovation) SQLParamsRequest() []string {
	return []string{
		"author_id",
		"title",
		"tags",
		"problem",
		"descriptions",
		"effect",
		"state",
		"meta",
		"created_at",
		"updated_at",
		"deleted_at",
	}
}

type InnovationDetail struct {
	Innovation
	Author   *Profile                      `json:"author"`
	CoAuthor *[]*InnovationCoAuthorsDetail `json:"co_authors"`
	Expert   *Profile                      `json:"expert"`
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

type InnovationCoAuthorsDetail struct {
	InnovationCoAuthors
	Author *Profile `json:"author"`
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

type InnovationExpertsDetail struct {
	InnovationExperts
	Expert *Profile
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
