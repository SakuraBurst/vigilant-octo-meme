package models

import "time"

type Banner struct {
	ID        int                    `json:"id"`
	Tags      []int                  `json:"tag_ids"`
	Feature   int                    `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
