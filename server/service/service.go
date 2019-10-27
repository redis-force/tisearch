package service

import (
	"context"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/redis-force/tisearch/logging"
	"github.com/redis-force/tisearch/server/model"
	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	esHostsEnv, _ = os.LookupEnv("ES_URLS")
	dbDSNEnv, _   = os.LookupEnv("DB_DSN")
)

const (
	indexType            = "doc"
	tweetSuggestionIndex = "tisearch-tweet-suggestion"
	userSuggestionIndex  = "tisearch-user-suggestion"
)

type TiSearchService struct {
	esClient *elastic.Client
	dbClient *gorm.DB
}

func NewSearchService() (*TiSearchService, error) {
	esHosts := []string{"http://117.50.101.237:9200/"}
	if len(esHostsEnv) != 0 {
		esHosts = strings.Split(esHostsEnv, ",")
	}
	rawClient, err := elastic.NewClient(elastic.SetURL(esHosts...), elastic.SetSniff(false))
	if err != nil {
		logging.Warnf("create es client error %s", err)
		return nil, err
	}
	dbDSN := "root:@tcp(10.9.118.254:4000)/tisearch?charset=utf8&timeout=1s&parseTime=true"
	if len(dbDSNEnv) != 0 {
		dbDSN = dbDSNEnv
	}
	db, err := gorm.Open("mysql", dbDSN)
	db.BlockGlobalUpdate(false)
	if err != nil {
		return nil, err
	}
	db = db.Debug()
	s := &TiSearchService{
		esClient: rawClient,
		dbClient: db,
	}
	return s, nil
}

func (s *TiSearchService) SearchTweet(keyword string, query string) (results []model.Tweet, plans []model.SQLPlan, err error) {
	results = make([]model.Tweet, 0)
	sql := "SELECT /*+ SEARCH('" + keyword + "' IN NATURAL LANGUAGE MODE) */" + query
	if err = s.dbClient.Raw(sql).Scan(&results).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = nil
		}
	}
	if len(results) > 200 {
		results = results[:200]
	}
	plans = make([]model.SQLPlan, 0)
	s.dbClient.Raw("EXPLAIN " + sql).Scan(&plans)
	return
}

func (s *TiSearchService) SuggestTweet(keyword string) ([]string, error) {
	suggester := elastic.NewCompletionSuggester("tweet-suggestion").Text(keyword).Field("words")
	searchSource := elastic.NewSearchSource().
		Suggester(suggester).
		FetchSource(false).
		TrackScores(true)
	searchResult, err := s.esClient.Search().Index(tweetSuggestionIndex).Type("words").SearchSource(searchSource).Do(context.TODO())
	if err != nil {
		return nil, err
	}
	stickerSuggest := searchResult.Suggest["tweet-suggestion"]
	fmt.Println(searchResult.Suggest["tweet-suggestion"])
	var results []string
	for _, options := range stickerSuggest {
		for _, option := range options.Options {
			results = append(results, option.Text)
		}
	}
	return results, nil
}

func (s *TiSearchService) SearchUser(keyword string, query string) (results []model.User, plans []model.SQLPlan, err error) {
	results = make([]model.User, 0)
	sql := "SELECT /*+ SEARCH('" + keyword + "' IN NATURAL LANGUAGE MODE) */" + query
	if err := s.dbClient.Raw(sql).Scan(&results).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			err = nil
		}
	}
	if len(results) > 200 {
		results = results[:200]
	}
	plans = make([]model.SQLPlan, 0)
	s.dbClient.Raw("EXPLAIN " + sql).Scan(&plans)
	return
}

func (s *TiSearchService) Execute(query string) (affeted int64, plans []model.SQLPlan, err error) {
	affeted = s.dbClient.Exec(query).RowsAffected
	err = s.dbClient.Error
	plans = make([]model.SQLPlan, 0)
	s.dbClient.Raw("EXPLAIN " + query).Scan(&plans)
	return
}

func (s *TiSearchService) SuggestUser(keyword string) ([]string, error) {
	suggester := elastic.NewCompletionSuggester("user-suggestion").Text(keyword).Field("words")
	searchSource := elastic.NewSearchSource().
		Suggester(suggester).
		FetchSource(false).
		TrackScores(true)
	searchResult, err := s.esClient.Search().Index(userSuggestionIndex).Type("words").SearchSource(searchSource).Do(context.TODO())
	if err != nil {
		return nil, err
	}
	stickerSuggest := searchResult.Suggest["user-suggestion"]
	fmt.Println(searchResult.Suggest["user-suggestion"])
	var results []string
	for _, options := range stickerSuggest {
		for _, option := range options.Options {
			results = append(results, option.Text)
		}
	}
	return results, nil
}
