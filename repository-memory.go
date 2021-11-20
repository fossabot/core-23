package main

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

var _ Repository = &MemoryRepository{}

type MemoryRepository struct {
	items []Item
	mu    sync.Mutex
}

func (repo *MemoryRepository) Insert(_ context.Context, item Item) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i := range repo.items {
		if repo.items[i].UUID() == item.UUID() {
			return errors.Errorf("item with uuid '%s' already exists", item.UUID())
		}

		if repo.items[i].Name() == item.Name() {
			return errors.Errorf("item with name '%s' already exists", item.Name())
		}
	}

	repo.items = append(repo.items, item)

	return nil
}

func (repo *MemoryRepository) List(_ context.Context) ([]Item, error) {
	return repo.items, nil
}

func (repo *MemoryRepository) FindByName(_ context.Context, name string) (Item, error) {
	for i := range repo.items {
		if repo.items[i].Name() == name {
			return repo.items[i], nil
		}
	}

	return nil, ErrItemNotFound
}

func (repo *MemoryRepository) Replace(_ context.Context, itemUUID string, item Item) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i := range repo.items {
		if repo.items[i].UUID() == itemUUID {
			repo.items[i] = item

			return nil
		}
	}

	return errors.Errorf("item with uuid '%s' doesn't exist", itemUUID)
}

func (repo *MemoryRepository) Delete(_ context.Context, itemUUID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i := range repo.items {
		if repo.items[i].UUID() == itemUUID {
			repo.items = append(repo.items[:i], repo.items[i+1:]...)

			return nil
		}
	}

	return errors.Errorf("item with uuid '%s' doesn't exist", itemUUID)
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make([]Item, 0),
		mu:    sync.Mutex{},
	}
}
