package main

import (
	"context"
	"sync"

	"github.com/pkg/errors"
)

var _ Repository = &MemoryRepository{}

type MemoryRepository struct {
	items []map[string]interface{}
	mu    sync.Mutex
}

func (repo *MemoryRepository) Insert(_ context.Context, item map[string]interface{}) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i := range repo.items {
		if repo.items[i][ItemFieldUUID] == item[ItemFieldUUID] {
			return errors.Errorf("item with uuid '%s' already exists", item[ItemFieldUUID])
		}
	}

	repo.items = append(repo.items, item)

	return nil
}

func (repo *MemoryRepository) List(_ context.Context) ([]map[string]interface{}, error) {
	return repo.items, nil
}

func (repo *MemoryRepository) Find(_ context.Context, itemUUID string) (map[string]interface{}, error) {
	for i := range repo.items {
		if repo.items[i][ItemFieldUUID] == itemUUID {
			return repo.items[i], nil
		}
	}

	return nil, nil
}

func (repo *MemoryRepository) Replace(_ context.Context, itemUUID string, item map[string]interface{}) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	for i := range repo.items {
		if repo.items[i][ItemFieldUUID] == itemUUID {
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
		if repo.items[i][ItemFieldUUID] == itemUUID {
			repo.items = append(repo.items[:i], repo.items[i+1:]...)

			return nil
		}
	}

	return errors.Errorf("item with uuid '%s' doesn't exist", itemUUID)
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: make([]map[string]interface{}, 0),
		mu:    sync.Mutex{},
	}
}
