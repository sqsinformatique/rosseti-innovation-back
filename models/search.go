package models

type Search struct {
	Value     string                 `json:"value"`
	ExtFilter map[string]interface{} `json:"ext_filter"`
}
