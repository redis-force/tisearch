package storage

import (
	"context"

	"github.com/redis-force/tisearch/model"
)

type Store interface {
	Create(ctx context.Context, db, table string, fields []model.Field) error
	Put(ctx context.Context, db string, table string, docID int64, fields []model.Field) error
	Search(ctx context.Context, db string, table string, query string) ([]model.SearchResult, error)
}

var defaultStore Store

func Create(ctx context.Context, db, table string, fields []model.Field) error {
	return defaultStore.Create(ctx, db, table, fields)
}

func Put(ctx context.Context, db, table string, docID int64, fields []model.Field) error {
	return defaultStore.Put(ctx, db, table, docID, fields)
}

func Search(ctx context.Context, db, table string, query string) ([]model.SearchResult, error) {
	return defaultStore.Search(ctx, db, table, query)
}
