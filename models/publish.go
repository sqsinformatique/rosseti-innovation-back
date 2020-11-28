package models

type Publish struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
	Type    string `json:"type"`
}
