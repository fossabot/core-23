package repository

import (
	"context"

	"github.com/nasermirzaei89/core/internal/core"
	"github.com/pkg/errors"
)

type Item interface {
	Insert(ctx context.Context, item core.Item) (err error)
	ListByType(ctx context.Context, typ string) (items []core.Item, err error)
	FindByTypeAndName(ctx context.Context, typ, name string) (item core.Item, err error)
	Replace(ctx context.Context, itemUUID string, item core.Item) (err error)
	Delete(ctx context.Context, itemUUID string) (err error)
}

var ErrItemNotFound = errors.New("item not found")
