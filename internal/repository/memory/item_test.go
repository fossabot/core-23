package memory_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nasermirzaei89/core/internal/core"
	"github.com/nasermirzaei89/core/internal/repository"
	"github.com/nasermirzaei89/core/internal/repository/memory"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestItemRepository_Insert(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("Insert 1", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		typ := "baz"

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      typ,
			Name:      "foo",
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item)
		assert.NoError(t, err)

		res, err := itemRepo.ListByType(ctx, typ)
		require.NoError(t, err)

		assert.Len(t, res, 1)

		assert.EqualValues(t, item, res[0])
	})

	t.Run("Insert 2", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		typ := "baz"

		items := []core.Item{
			{
				UUID:      uuid.NewString(),
				Type:      typ,
				Name:      "foo",
				Data:      nil,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				UUID:      uuid.NewString(),
				Type:      typ,
				Name:      "foo2",
				Data:      nil,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		for i := range items {
			err := itemRepo.Insert(ctx, items[i])
			assert.NoError(t, err)
		}

		res, err := itemRepo.ListByType(ctx, typ)
		require.NoError(t, err)

		assert.EqualValues(t, items, res)
	})

	t.Run("Duplicate UUID", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		typ := "baz"

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      typ,
			Name:      "foo",
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item)
		assert.NoError(t, err)

		err = itemRepo.Insert(ctx, item)
		assert.Error(t, err)
	})

	t.Run("Duplicate Type+Name", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		typ := "baz"

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      typ,
			Name:      "foo",
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item)
		assert.NoError(t, err)

		item.UUID = uuid.NewString()

		err = itemRepo.Insert(ctx, item)
		assert.Error(t, err)
	})
}

func TestItemRepository_FindByTypeAndName(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("Found", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		name := "foo"
		typ := "bar"

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      typ,
			Name:      name,
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item)
		require.NoError(t, err)

		res, err := itemRepo.GetByTypeAndName(ctx, typ, name)
		assert.NoError(t, err)

		assert.EqualValues(t, item, *res)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		name := "fee"
		typ := "bar"

		_, err := itemRepo.GetByTypeAndName(ctx, name, typ)
		assert.True(t, errors.Is(err, repository.ErrItemNotFound))
	})
}

func TestItemRepository_Replace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("Replace", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		name := "foo"
		typ := "bar"

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      typ,
			Name:      name,
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item)
		require.NoError(t, err)

		item.Data = map[string]interface{}{
			"foo": "bar",
		}

		err = itemRepo.Replace(ctx, item.UUID, item)
		assert.NoError(t, err)

		res, err := itemRepo.GetByTypeAndName(ctx, typ, name)
		require.NoError(t, err)

		assert.EqualValues(t, item, *res)
	})

	t.Run("Change UUID", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		name := "foo"
		typ := "bar"

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      typ,
			Name:      name,
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item)
		require.NoError(t, err)

		oldUUID := item.UUID

		item.UUID = uuid.NewString()

		err = itemRepo.Replace(ctx, oldUUID, item)
		assert.Error(t, err)
	})

	t.Run("Conflict on Type and Name", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		item1 := core.Item{
			UUID:      uuid.NewString(),
			Type:      "bar",
			Name:      "foo",
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item1)
		require.NoError(t, err)

		item2 := core.Item{
			UUID:      uuid.NewString(),
			Type:      "bar",
			Name:      "fee",
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = itemRepo.Insert(ctx, item2)
		require.NoError(t, err)

		item2.Name = item1.Name

		err = itemRepo.Replace(ctx, item2.UUID, item2)
		assert.Error(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      "foo",
			Name:      "bar",
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Replace(ctx, item.UUID, item)
		assert.Error(t, err)
	})
}

func TestItemRepository_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("Normal", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		item := core.Item{
			UUID:      uuid.NewString(),
			Type:      "foo",
			Name:      "bar",
			Data:      nil,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := itemRepo.Insert(ctx, item)
		require.NoError(t, err)

		err = itemRepo.Delete(ctx, item.UUID)
		assert.NoError(t, err)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		itemRepo := memory.NewItemRepository()

		err := itemRepo.Delete(ctx, uuid.NewString())
		assert.Error(t, err)
	})
}
