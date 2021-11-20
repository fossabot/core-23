package main

type Item map[string]interface{}

func (item Item) UUID() string {
	return item[ItemFieldUUID].(string)
}

func (item Item) Type() string {
	return item[ItemFieldType].(string)
}

func (item Item) Name() string {
	return item[ItemFieldName].(string)
}

const (
	ItemFieldUUID      = "uuid"
	ItemFieldType      = "type"
	ItemFieldName      = "name"
	ItemFieldCreatedAt = "createdAt"
	ItemFieldUpdatedAt = "updatedAt"
)

const (
	NameRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
	TypeRegex = "^[a-z][a-z0-9-]{1,254}[a-z0-9]$"
)
