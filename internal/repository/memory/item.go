package memory

import (
	"context"
	"sync"

	"github.com/nasermirzaei89/core/internal/core"
	"github.com/nasermirzaei89/core/internal/repository"
	"github.com/pkg/errors"
)

var _ repository.Item = &ItemRepository{}

type ItemRepository struct {
	items []core.Item
	mu    sync.Mutex
}

func (repo *ItemRepository) Insert(_ context.Context, item core.Item) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i := range repo.items {
		if repo.items[i].UUID == item.UUID {
			return errors.Errorf("item with uuid '%s' already exists", item.UUID)
		}

		if repo.items[i].Type == item.Type && repo.items[i].Name == item.Name {
			return errors.Errorf("item with type '%s' and name '%s' already exists", item.Type, item.Name)
		}
	}

	repo.items = append(repo.items, item)

	return nil
}

func (repo *ItemRepository) ListByType(_ context.Context, typ string) ([]core.Item, error) {
	res := make([]core.Item, 0)

	for i := range repo.items {
		if repo.items[i].Type == typ {
			res = append(res, repo.items[i])
		}
	}

	return res, nil
}

func (repo *ItemRepository) FindByTypeAndName(_ context.Context, typ, name string) (core.Item, error) {
	for i := range repo.items {
		if repo.items[i].Type == typ && repo.items[i].Name == name {
			return repo.items[i], nil
		}
	}

	return core.Item{}, repository.ErrItemNotFound
}

func (repo *ItemRepository) Replace(_ context.Context, itemUUID string, item core.Item) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if item.UUID != itemUUID {
		return errors.New("field uuid is immutable")
	}

	for i := range repo.items {
		if repo.items[i].Type == item.Type && repo.items[i].Name == item.Name && repo.items[i].UUID != itemUUID {
			return errors.Errorf("item with type '%s' and name '%s' already exists", item.Type, item.Name)
		}
	}

	for i := range repo.items {
		if repo.items[i].UUID == itemUUID {
			repo.items[i] = item

			return nil
		}
	}

	return errors.Errorf("item with uuid '%s' doesn't exist", itemUUID)
}

func (repo *ItemRepository) Delete(_ context.Context, itemUUID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i := range repo.items {
		if repo.items[i].UUID == itemUUID {
			repo.items = append(repo.items[:i], repo.items[i+1:]...)

			return nil
		}
	}

	return errors.Errorf("item with uuid '%s' doesn't exist", itemUUID)
}

func NewMemoryRepository() *ItemRepository {
	return &ItemRepository{
		items: make([]core.Item, 0),
		mu:    sync.Mutex{},
	}
}
