package model

import (
	"encoding/json"
	"strings"
	"time"
)

type Tweet struct {
	ID      int64  `json:"id" gorm:"column:id"`
	User    string `json:"user" gorm:"column:user"`
	Content string `json:"content" gorm:"column:content"`
	Time    string `json:"time" gorm:"column:time"`
}
type User struct {
	ID          int64     `json:"id" gorm:"primary_key,column:id"`
	Gender      int       `json:"gender" gorm:"column:gender"`
	Name        string    `json:"name" gorm:"column:name"`
	Location    string    `json:"location" gorm:"column:location"`
	Picture     string    `json:"picture" gorm:"column:picture"`
	Labels      string    `json:"-" gorm:"column:labels"`
	Interests   []string  `json:"interests" gorm:"-"`
	Cars        []string  `json:"cars" gorm:"-"`
	Birthday    time.Time `json:"birthday" gorm:"column:birthday"`
	Coordinates string    `json:"coordinates" gorm:"column:coordinates"`
}
type user struct {
	ID          int64    `json:"id" `
	Gender      string   `json:"gender"`
	Name        string   `json:"name"`
	Location    string   `json:"location"`
	Picture     string   `json:"picture"`
	Interests   []string `json:"interests"`
	Cars        []string `json:"cars"`
	Birthday    string   `json:"birthday"`
	Coordinates []string `json:"coordinates"`
}

type Suggestion struct {
	Suggestion []string `json:"suggestion"`
}

type SuggestionResponse struct {
	Code  int        `json:"code"`
	Error string     `json:"error"`
	Data  Suggestion `json:"data"`
}

type SearchResponse struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Type  string      `json:"type"`
	Data  interface{} `json:"data"`
	Plans []SQLPlan   `json:"plans"`
}
type SQLPlan struct {
	ID           string `json:"id" gorm:"column:id"`
	Count        string `json:"count" gorm:"column:count"`
	Task         string `json:"task" gorm:"column:task"`
	OperatorInfo string `json:"operation_info" gorm:"column:operator info"`
}

func (u User) MarshalJSON() ([]byte, error) {
	us := user{
		ID:          u.ID,
		Gender:      "male",
		Name:        u.Name,
		Location:    u.Location,
		Picture:     u.Picture,
		Birthday:    u.Birthday.Format("2006-01-02 15:04:05"),
		Coordinates: strings.Split(u.Coordinates, ","),
	}
	if u.Gender == 2 {
		us.Gender = "female"
	}
	labels := strings.Split(u.Labels, " ")
	var (
		cars      []string
		interests []string
	)
	for _, l := range labels {
		if strings.HasPrefix(l, "car_") {
			cars = append(cars, strings.TrimPrefix(l, "car_"))
		} else if strings.HasPrefix(l, "interest_") {
			interests = append(interests, strings.TrimPrefix(l, "interest_"))
		}
	}
	us.Interests = interests
	us.Cars = cars
	return json.Marshal(us)
}
