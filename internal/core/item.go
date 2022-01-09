package core

import "time"

const (
	NameRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
	TypeRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
)

type Item struct {
	UUID      string                 `json:"uuid"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Data      map[string]interface{} `json:",inline"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
}

type CreateItemRequest struct {
	Name string                 `json:"name"`
	Data map[string]interface{} `json:",inline"`
}

type ReplaceItemRequest struct {
	Data map[string]interface{} `json:",inline"`
}

type ListItemsResponse struct {
	Items []Item `json:"items"`
}
