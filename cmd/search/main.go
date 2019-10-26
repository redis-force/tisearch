package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/redis-force/tisearch/storage"
)

var (
	table = flag.String("table", "tweets", "table name to search")
)

func search(query string) {
	result, err := storage.Search(context.Background(), "tisearch", *table, query)
	if err != nil {
		panic(err)
	}
	lines := make([]string, 0, 200)
	maxLen := 0
	for _, r := range result.Rows {
		line := fmt.Sprintf("| %d ", r.DocID)
		if len(r.Render) == 1 {
			for _, v := range r.Render {
				line += fmt.Sprintf("| %s ", v[0])
			}
		} else {
			// verbose mode
			for c, v := range r.Render {
				line += fmt.Sprintf("| %s:%s ", c, v[0])
			}
		}
		lines = append(lines, line)
		l := len(line) + 1
		if l > maxLen {
			maxLen = l
		}
	}
	fmt.Println(strings.Repeat("=", maxLen))
	for _, line := range lines {
		fmt.Printf("%s%s|\n", line, strings.Repeat(" ", maxLen-len(line)-1))
	}
	fmt.Println(strings.Repeat("=", maxLen))
}

func main() {
	flag.Parse()
	for _, query := range os.Args[1:] {
		search(query)
	}
}
