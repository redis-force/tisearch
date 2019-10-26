package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/redis-force/tisearch/logging"
	"github.com/redis-force/tisearch/model"
	"github.com/redis-force/tisearch/storage/es"
)

func main() {
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
	file, err := os.Open("training.1600000.processed.noemoticon.csv")
	if err != nil {
		logging.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cnt := int64(1)
	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), ",", 6)
		fields = []model.Field{
			model.Field{
				Name:  "user",
				Value: strings.Trim(line[4], "\""),
			},
			model.Field{
				Name:  "content",
				Value: strings.Trim(line[5], "\""),
			},
		}
		err = es.Put(context.TODO(), "twitter", "twitter", cnt, fields)
		cnt++
		if err != nil {
			logging.Fatal(err)
		}

	}

	if err := scanner.Err(); err != nil {
		logging.Fatal(err)
	}
}
