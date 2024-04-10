package models

type Banner struct {
	Tags     []int                  `json:"tag_ids"`
	Feature  int                    `json:"feature_id"`
	Content  map[string]interface{} `json:"content"`
	IsActive bool                   `json:"is_active"`
}