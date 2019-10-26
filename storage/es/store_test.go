package es

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/redis-force/tisearch/model"
)

func TestEsStore(t *testing.T) {
	es, err := NewEsStore()
	if err != nil {
		t.Fatal(err)
	}
	fields := []model.Field{
		model.Field{
			Name:  "message",
			Value: "this is message2 content a ha",
		},
		model.Field{
			Name:  "user",
			Value: "this is user_name",
		},
	}
	err = es.Create(context.TODO(), "test", "test", fields)
	if err != nil {
		t.Fatal(err)
	}
	err = es.Put(context.TODO(), "test", "test", 1, fields)
	if err != nil {
		t.Fatal(err)
	}
	result, err := es.Search(context.TODO(), "test", "test", `"this is message2"`)
	if err != nil {
		t.Fatal(err)
	}
	err = es.Delete(context.TODO(), "tisearch", "tweets2", 1)
	if err != nil {
		t.Fatal(err)
	}
	err = es.Delete(context.TODO(), "tisearch", "tweets2", 2)
	if err != nil {
		t.Fatal(err)
	}
	err = es.Delete(context.TODO(), "tisearch", "tweets2", 3)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}

func TestTwitterLoad(t *testing.T) {
	es, err := NewEsStore()
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}
	file, err := os.Open("/Users/melody/Downloads/trainingandtestdata/training.1600000.processed.noemoticon.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cnt := int64(1)
	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), ",", 6)
		fields = []model.Field{
			model.Field{
				Name:  "user",
				Value: line[4],
			},
			model.Field{
				Name:  "content",
				Value: line[5],
			},
		}
		err = es.Put(context.TODO(), "twitter", "twitter", cnt, fields)
		cnt++
		if err != nil {
			t.Fatal(err)
		}

	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}
