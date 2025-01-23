package v1

import (
	"context"
)

type inMemoryRepository struct{}

func newInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{}
}

func (repo *inMemoryRepository) getHelloWorldMessage(_ context.Context) (string, error) {
	return "Hello World", nil
}
