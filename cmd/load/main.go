package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis-force/tisearch/logging"
	"github.com/redis-force/tisearch/model"
	"github.com/redis-force/tisearch/storage/es"
)

type row struct {
	polarity int
	id       int64
	time     time.Time
	user     string
	content  string
}

var (
	mode  = flag.String("mode", "es", "load mode: es/tidb")
	input = flag.String("input-file", "training.1600000.processed.noemoticon.csv", "load input file")
	tidb  = flag.String("tidb", "root:@tcp(127.0.0.1:4000)/tisearch?autocommit=false", "tidb connection string")

	inputChannel = make(chan row)
	wg           = sync.WaitGroup{}
)

func loadEs() {
	defer wg.Done()
	es, err := es.NewEsStore()
	if err != nil {
		logging.Fatal(err)
	}
	fields := []model.Field{
		model.Field{
			Name: "user",
		},
		model.Field{
			Name: "content",
		},
	}
	err = es.Create(context.TODO(), "twitter", "twitter", fields)
	if err != nil {
		logging.Fatal(err)
	}
	cnt := int64(1)
	for r := range inputChannel {
		fields = []model.Field{
			model.Field{
				Name:  "user",
				Value: r.user,
			},
			model.Field{
				Name:  "content",
				Value: r.content,
			},
		}
		err = es.Put(context.TODO(), "twitter", "twitter", cnt, fields)
		cnt++
		if err != nil {
			logging.Fatal(err)
		}
	}
}

func loadTiDB() {
	defer wg.Done()
	db, err := sql.Open("mysql", *tidb)
	if err != nil {
		logging.Fatal(err)
	}
	defer db.Close()

	var pstmt *sql.Stmt
	var tx *sql.Tx

	cnt := 0
	for r := range inputChannel {
		if tx == nil {
			tx, err = db.Begin()
			if err != nil {
				logging.Error(err)
				continue
			}
		}
		if pstmt == nil {
			pstmt, err = tx.Prepare("INSERT INTO tisearch.tweets(id, time, user, content, polarity) VALUES(?, ?, ?, ?, ?)")
			if err != nil {
				logging.Fatal(err)
				continue
			}
		}
		_, err := pstmt.Exec(r.id, r.time, r.user, r.content, r.polarity)
		if err != nil {
			logging.Error(err)
			continue
		}
		cnt++
		if cnt%1000 == 0 {
			tx.Commit()
			pstmt.Close()
			tx = nil
			pstmt = nil
			logging.Infof("Committed %d rows of data", cnt)
		}
	}
}

func main() {
	flag.Parse()
	file, err := os.Open(*input)
	if err != nil {
		logging.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	wg.Add(1)
	switch *mode {
	case "tidb":
		go loadTiDB()
	default:
		go loadEs()
	}

	scanned := 0
	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), ",", 6)
		for i, p := range line {
			line[i] = strings.Trim(p, "\"")
		}
		polarity, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			logging.Errorf("caught error: %v", err)
			continue
		}
		id, err := strconv.ParseInt(line[1], 10, 64)
		if err != nil {
			logging.Errorf("caught error: %v", err)
			continue
		}
		timestamp, err := time.Parse("Mon Jan 02 15:04:05 MST 2006", line[2])
		if err != nil {
			logging.Errorf("caught error: %v", err)
			continue
		}
		user := line[4]
		content := line[5]
		inputChannel <- row{
			polarity: int(polarity),
			id:       id,
			time:     timestamp,
			user:     user,
			content:  content,
		}
		scanned++
		if scanned%1000 == 0 {
			logging.Infof("scanned %v rows", scanned)
		}
	}

	if err := scanner.Err(); err != nil {
		logging.Fatal(err)
	}

	wg.Wait()
}
