package core

import "time"

type Item struct {
	UUID      string
	Type      string
	Name      string
	Data      map[string]interface{}
	CreatedAt time.Time
	UpdatedAt time.Time
}

const (
	NameRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
	TypeRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
)
