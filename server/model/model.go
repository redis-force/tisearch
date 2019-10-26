package model

import "time"

type Tweet struct {
	Creator    string    `json:"creator" gorm:"column:creator"`
	Content    string    `json:"content" gorm:"column:content"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

type Album struct {
	Gender      string    `json:"gender"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Picture     string    `json:"picture"`
	Cars        string    `json:"cars"`
	Interests   string    `json:"interests"`
	Birthday    time.Time `json:"birthday"`
	Coordinates string    `json:"coordinates"`
}

type Suggestion struct {
	Suggestion []string `json:"suggestion"`
}

type TweetSuggestionResponse struct {
	Code  int        `json:"code"`
	Error string     `json:"error"`
	Data  Suggestion `json:"data"`
}

type AlbumSearchResponse struct {
	Code  int     `json:"code"`
	Error string  `json:"error"`
	Data  []Tweet `json:"data"`
}

type TweetSearchResponse struct {
	Code  int     `json:"code"`
	Error string  `json:"error"`
	Data  []Tweet `json:"data"`
}

type AlbumSuggestionResponse struct {
	Code  int        `json:"code"`
	Error string     `json:"error"`
	Data  Suggestion `json:"data"`
}
