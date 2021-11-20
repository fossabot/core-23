package main

import (
	"context"

	"github.com/pkg/errors"
)

type Repository interface {
	Insert(ctx context.Context, item map[string]interface{}) (err error)
	List(ctx context.Context) (items []map[string]interface{}, err error)
	FindByName(ctx context.Context, name string) (item map[string]interface{}, err error)
	Replace(ctx context.Context, itemUUID string, item map[string]interface{}) (err error)
	Delete(ctx context.Context, itemUUID string) (err error)
}

var ErrItemNotFound = errors.New("item not found")
