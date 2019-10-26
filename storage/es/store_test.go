package es

import (
	"context"
	"fmt"
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
			Value: "this is message content a ha",
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
	result, err := es.Search(context.TODO(), "test", "test", "user:*name")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
