package main

type Item map[string]interface{}

func (item Item) UUID() string {
	return item[ItemFieldUUID].(string)
}

func (item Item) Name() string {
	return item[ItemFieldUUID].(string)
}

const (
	ItemFieldUUID      = "uuid"
	ItemFieldName      = "name"
	ItemFieldCreatedAt = "createdAt"
	ItemFieldUpdatedAt = "updatedAt"
)

const (
	NameRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
)
