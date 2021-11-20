package main

import (
	"context"

	"github.com/pkg/errors"
)

type Repository interface {
	Insert(ctx context.Context, item Item) (err error)
	ListByType(ctx context.Context, typ string) (items []Item, err error)
	FindByTypeAndName(ctx context.Context, typ, name string) (item Item, err error)
	Replace(ctx context.Context, itemUUID string, item Item) (err error)
	Delete(ctx context.Context, itemUUID string) (err error)
}

var ErrItemNotFound = errors.New("item not found")
