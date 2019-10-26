package es

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/redis-force/tisearch/logging"
	"github.com/redis-force/tisearch/model"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	esHostsEnv, _ = os.LookupEnv("ES_URLS")
)

const (
	indexType = "doc"
)

type EsStore struct {
	esClient *elastic.Client
	ctx      context.Context
}

func NewEsStore() (*EsStore, error) {
	s := &EsStore{
		ctx: context.Background(),
	}
	esHosts := []string{"http://117.50.101.237:9200/"}
	if len(esHostsEnv) != 0 {
		esHosts = strings.Split(esHostsEnv, ",")
	}
	rawClient, err := elastic.NewClient(elastic.SetURL(esHosts...), elastic.SetSniff(false))
	if err != nil {
		logging.Warnf("create es client error %s", err)
		return nil, err
	}
	s.esClient = rawClient
	return s, nil
}

func (s *EsStore) Create(ctx context.Context, db, table string, fields []model.Field) (err error) {
	index := indexName(db, table)
	exists, _ := s.esClient.IndexExists(index).Do(ctx)
	if exists {
		logging.Warnf("create es mapping (db=%s, table=%s), index already exists %s", db, table, err)
		return nil
	}

	properties := make(map[string]interface{})
	for _, f := range fields {
		properties[f.Name] = map[string]interface{}{"type": "string"}
	}
	propertiesStr, _ := json.Marshal(map[string]interface{}{"properties": properties})
	// _, err = s.esClient.PutMapping().Type(indexType).Index(index).BodyString(string(propertiesStr)).Do(ctx)
	body := fmt.Sprintf(mapping, indexType, propertiesStr)
	s.esClient.CreateIndex(index).BodyString(body).Do(ctx)
	logging.Debugf("create db=%s table=%s mapping=%q\n", db, table, propertiesStr)
	return
}

func (s *EsStore) Put(ctx context.Context, db string, table string, docID int64, fields []model.Field) error {
	index := indexName(db, table)
	data := make(map[string]interface{})
	for _, f := range fields {
		data[f.Name] = f.Value
	}
	_, err := s.esClient.Index().Index(index).Type(indexType).Id(strconv.Itoa(int(docID))).BodyJson(data).Do(ctx)
	logging.Debugf("write index to %s, body %v", index, data)
	if err != nil {
		return err
	}
	return nil
}

func (s *EsStore) Search(ctx context.Context, db string, table string, query string) (*model.SearchResult, error) {
	index := indexName(db, table)
	if exists, err := s.esClient.IndexExists(index).Do(ctx); !exists || err != nil {
		return nil, err
	}
	indexMapping, err := s.esClient.GetMapping().Index(index).Type(indexType).Do(ctx)
	if err != nil {
		return nil, err
	}
	mm := indexMapping[index].(map[string]interface{})["mappings"].(map[string]interface{})[indexType].(map[string]interface{})["properties"].(map[string]interface{})
	highlighter := elastic.NewHighlight().PreTags("<hit>").PostTags("</hit>")
	for k := range mm {
		highlighter.Field(k)
	}
	q := elastic.NewQueryStringQuery(query)
	searchResult, err := s.esClient.Search(index).Query(q).Type(indexType).Size(1000).IgnoreUnavailable(true).Highlight(highlighter).Do(ctx)
	if err != nil {
		logging.Warnf("search db=%s table %s query=%s error", db, table, query)
		return nil, err
	}
	if searchResult.Hits == nil || len(searchResult.Hits.Hits) == 0 {
		return nil, nil
	}
	result := new(model.SearchResult)
	for _, hit := range searchResult.Hits.Hits {
		_ = hit
		docID, _ := strconv.Atoi(hit.Id)
		if docID != 0 {
			var row model.Row
			row.DocID = int64(docID)
			row.Render = hit.Highlight
			result.Rows = append(result.Rows, row)
		}
	}
	return result, nil
}
