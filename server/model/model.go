package model

import (
	"time"
)

type Tweet struct {
	ID      int64     `json:"id" gorm:"column:id"`
	User    string    `json:"user" gorm:"column:user"`
	Content string    `json:"content" gorm:"column:content"`
	Time    time.Time `json:"time" gorm:"column:time,type:timestamp"`
}

type User struct {
	ID       int64  `json:"id" gorm:"primary_key,column:id"`
	Gender   int    `json:"gender" gorm:"column:gender"`
	Name     string `json:"name" gorm:"column:name"`
	Location string `json:"location" gorm:"column:location"`
	Picture  string `json:"picture" gorm:"column:picture"`
	Labels   string `json:"cars" gorm:"column:labels"`
	// Interests   string    `json:"interests" gorm:"column:interests"`
	Birthday    time.Time `json:"birthday" gorm:"column:birthday"`
	Coordinates string    `json:"coordinates" gorm:"column:coordinates"`
}

type Suggestion struct {
	Suggestion []string `json:"suggestion"`
}

type TweetSuggestionResponse struct {
	Code  int        `json:"code"`
	Error string     `json:"error"`
	Data  Suggestion `json:"data"`
}

type UserSearchResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
	Data  []User `json:"data"`
}

type TweetSearchResponse struct {
	Code  int     `json:"code"`
	Error string  `json:"error"`
	Data  []Tweet `json:"data"`
}

type UserSuggestionResponse struct {
	Code  int        `json:"code"`
	Error string     `json:"error"`
	Data  Suggestion `json:"data"`
}
