package service

import (
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
	indexType  = "doc"
	tweetIndex = "tisearch-tweet"
	albumIndex = "tisearch-album"
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
	dbDSN := "root:@tcp(10.9.118.254:3306)/tisearch?charset=utf8&timeout=1s"
	if len(dbDSNEnv) != 0 {
		dbDSN = dbDSNEnv
	}
	db, err := gorm.Open("mysql", dbDSN)
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

// SELECT /*+ SEARCH(‘run out of time’ IN NATURAL LANGUAGE MODE) */ * from tweets;
// SELECT /*+ SEARCH(‘+run +out -money’ IN BOOLEAN MODE) */ * from tweets;
//
//
func (s *TiSearchService) SearchTweetByKeyword(keyword string) ([]model.Tweet, error) {
	results := make([]model.Tweet, 0)
	if err := s.dbClient.Raw("SELECT /*+ SEARCH('?' IN NATURAL LANGUAGE MODE) */ * from tweets;", keyword).Scan(&results).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
	}
	return results, nil
}

func (s *TiSearchService) SearchAlbum(keyword string) ([]model.Album, error) {
	results := make([]model.Album, 0)
	if err := s.dbClient.Raw("SELECT /*+ SEARCH('?' IN NATURAL LANGUAGE MODE) */ * from tweets;", keyword).Scan(&results).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
	}
	return results, nil
}
