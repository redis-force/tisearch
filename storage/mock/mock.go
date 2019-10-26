package mock

import (
	"context"

	"github.com/redis-force/tisearch/model"
)

type MockStore struct {
}

func (s *MockStore) Create(ctx context.Context, db, table string, fields []model.Field) error {
}

func (s *MockStore) Put(ctx context.Context, db string, table string, docID int64, fields []model.Field) error {
}

func (s *MockStore) Search(ctx context.Context, db string, table string, query string) (*model.SearchResult, error) {
}

func (s *MockStore) Delete(ctx context.Context, db string, table string, docID int64) (*model.SearchResult, error) {
}
