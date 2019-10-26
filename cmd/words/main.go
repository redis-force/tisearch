package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/redis-force/tisearch/logging"
)

var (
	input       = flag.String("input-file", "training.1600000.processed.noemoticon.csv", "load input file")
	wordsOutput = flag.String("words", "words.csv", "words output file")
	usersOutput = flag.String("users", "users.csv", "users output file")
)

func writeTo(filename string, data map[string]int) {
	if len(data) > 0 {
		out, err := os.Create(filename)
		if err != nil {
			logging.Fatal(err)
		} else {
			for k, c := range data {
				out.WriteString(fmt.Sprintf("%d, %s\n", c, k))
			}
			out.Close()
		}
	}
}

func tidy(s string) string {
	return strings.TrimSpace(strings.Trim(strings.Trim(strings.Trim(strings.Trim(strings.Trim(s, "."), "#"), "!"), "("), ")"))
}

func main() {
	flag.Parse()
	file, err := os.Open(*input)
	if err != nil {
		logging.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	users := make(map[string]int)
	words := make(map[string]int)

	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), ",", 6)
		user := strings.Trim(line[4], "\"")
		content := strings.ReplaceAll(strings.ToLower(strings.Trim(line[5], "\"")), ",", " ")
		if c, ok := users[user]; ok {
			users[user] = c + 1
		} else {
			users[user] = 1
		}
		for _, word := range strings.Split(content, " ") {
			word = tidy(word)
			if strings.HasPrefix(word, "@") || strings.HasPrefix(word, "http://") || strings.HasPrefix(word, "https://") {
				continue
			}
			if c, ok := words[word]; ok {
				words[word] = c + 1
			} else {
				words[word] = 1
			}
		}
	}

	writeTo(*usersOutput, users)
	writeTo(*wordsOutput, words)

	if err := scanner.Err(); err != nil {
		logging.Fatal(err)
	}
}
