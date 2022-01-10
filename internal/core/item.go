package core

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

const (
	NameRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
	TypeRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
)

type Item struct {
	UUID      string `json:"uuid"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Data      map[string]interface{}
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (item Item) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})

	for k, v := range item.Data {
		m[k] = v
	}

	m["uuid"] = item.UUID
	m["type"] = item.Type
	m["name"] = item.Name
	m["createdAt"] = item.CreatedAt.Format(time.RFC3339)
	m["updatedAt"] = item.UpdatedAt.Format(time.RFC3339)

	res, err := json.Marshal(m)
	if err != nil {
		return nil, errors.Wrap(err, "error on marshal json")
	}

	return res, nil
}

func (item *Item) UnmarshalJSON(data []byte) error {
	var res map[string]interface{}

	err := json.Unmarshal(data, &res)
	if err != nil {
		return errors.Wrap(err, "error on unmarshal json to map")
	}

	item.Data = make(map[string]interface{})

	for k, v := range res {
		switch k {
		case "uuid":
			f, ok := v.(string)
			if !ok {
				return errors.New("field uuid is not string")
			}

			item.UUID = f
		case "type":
			f, ok := v.(string)
			if !ok {
				return errors.New("field type is not string")
			}

			item.Type = f
		case "name":
			f, ok := v.(string)
			if !ok {
				return errors.New("field name is not string")
			}

			item.Name = f
		case "createdAt":
			f, ok := v.(string)
			if !ok {
				return errors.New("field createdAt is not string")
			}

			t, err := time.Parse(time.RFC3339, f)
			if err != nil {
				return errors.Wrap(err, "error on parse createdAt time string")
			}

			item.CreatedAt = t
		case "updatedAt":
			f, ok := v.(string)
			if !ok {
				return errors.New("field updatedAt is not string")
			}

			t, err := time.Parse(time.RFC3339, f)
			if err != nil {
				return errors.Wrap(err, "error on parse updatedAt time string")
			}

			item.UpdatedAt = t
		default:
			item.Data[k] = v
		}
	}

	return nil
}

type ItemList struct {
	Items []Item `json:"items"`
}
