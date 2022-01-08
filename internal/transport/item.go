package transport

import (
	"time"

	"github.com/nasermirzaei89/core/internal/core"
)

type Item struct {
	UUID      string                 `json:"uuid"`
	Type      string                 `json:"type"`
	Name      string                 `json:"name"`
	Data      map[string]interface{} `json:",inline"`
	CreatedAt string                 `json:"createdAt"`
	UpdatedAt string                 `json:"updatedAt"`
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

func ItemFromEntity(item core.Item) Item {
	return Item{
		UUID:      item.UUID,
		Type:      item.Type,
		Name:      item.Name,
		Data:      item.Data,
		CreatedAt: item.CreatedAt.Format(time.RFC3339),
		UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
	}
}
