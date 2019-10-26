package es

import (
	"context"
	"os"
	"strings"

	"github.com/redis-force/tisearch/model"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	esHostsEnv, _ = os.LookupEnv("ES_URLS")
)

type EsStore struct {
	esClient *elastic.Client
	ctx      context.Context
}

func NewEsStore() (*EsStore, error) {
	s := &EsStore{
		ctx: context.Background(),
	}
	esHosts := []string{"http://10.9.120.175/"}
	if len(esHostsEnv) == 0 {
		esHosts = strings.Split(esHostsEnv, ",")
	}
	rawClient, err := elastic.NewClient(elastic.SetURL(esHosts...), elastic.SetSniff(true))
	if err != nil {
		return nil, err
	}
	s.esClient = rawClient
	return s, nil
}

func (s *EsStore) Create(ctx context.Context, db, table string, fields []model.Field) error {
	exists, _ := s.esClient.IndexExists(db + "-" + table).Do(ctx)
	if !exists {
		rsp, err := s.esClient.PutMapping().BodyJson(j).Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *EsStore) Put(ctx context.Context, db string, table string, docID int64, fields []model.Field) error {
	return nil
}

func (s *EsStore) Search(ctx context.Context, db string, table string, query string) ([]model.SearchResult, error) {
	return nil, nil
}
