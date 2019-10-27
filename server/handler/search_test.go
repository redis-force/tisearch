package handler

import (
	"fmt"
	"testing"
)

func TestMatch(t *testing.T) {
	// reg := regexp.MustCompile(regex)
	qs := []string{
		"SELECT id, time, tweets.user, polarity, content FROM tweets JOIN tweets_user on tweets.user = tweets_user.user WHERE MATCH(content) AGAINST ('mississipi')",
		"SELECT id, time, tweets.user, polarity, content FROM tweets JOIN tweets_user on tweets.user = tweets_user.user WHERE MATCH(content) AGAINST ('mississipi') AND time < '2009-05-18 00:00:00'",
		"SELECT id, name, location, picture, birthday, coordinates, gender, labels FROM users WHERE MATCH(labels) AGAINST ('interest_floorball car_alfa')",
		"SELECT id, name, location, picture, birthday, coordinates, gender, labels FROM users WHERE MATCH(labels) AGAINST ('interest:floorball car:alfa') AND birthday BETWEEN '1960-01-01 00:00:00' AND '1969-12-31 23:59:59'",
		"SELECT id, time, user, polarity, content FROM tweets WHERE MATCH(content) AGAINST ('清华大学')",
		"SELECT id, time, user, polarity, content FROM tweets WHERE MATCH(content) AGAINST ('清华大学')",
		"SELECT id, time, user, polarity, content FROM tweets WHERE MATCH(content) AGAINST ('清华大学')",
		"INSERT INTO tweets VALUES(10000, '2019/10/26 08:30:00', 'user', 2, '清华大学的前身清华学堂始建与1911年')",
		"DELETE FROM tweets WHERE id = 10000",
	}
	for _, q := range qs {
		fmt.Println(mathQuery(q))
	}
}
