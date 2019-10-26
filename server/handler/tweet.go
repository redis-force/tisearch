package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis-force/tisearch/server/model"
	"github.com/redis-force/tisearch/server/service"
)

type Handler struct {
	srv *service.TiSearchService
}

func NewHandler(srv *service.TiSearchService) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) SearchTweet(c *gin.Context) {
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		c.AbortWithStatusJSON(400, model.TweetSearchResponse{Code: 400, Error: "keyword 不能为空"})
		return
	}
	data, err := h.srv.SearchTweetByKeyword(keyword)
	if err != nil {
		c.AbortWithStatusJSON(500, model.TweetSearchResponse{Code: 500, Error: err.Error()})
		return
	}
	c.JSON(200, model.TweetSearchResponse{Code: 0, Error: "", Data: data})
}
func (h *Handler) SuggestTweet(c *gin.Context) {
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		c.AbortWithStatusJSON(400, model.TweetSuggestionResponse{Code: 400, Error: "keyword 不能为空"})
		return
	}
	data, err := h.srv.SuggestTweet(keyword)
	if err != nil {
		c.AbortWithStatusJSON(500, model.TweetSuggestionResponse{Code: 500, Error: err.Error()})
		return
	}
	sug := model.Suggestion{Suggestion: data}
	c.JSON(200, model.TweetSuggestionResponse{Code: 0, Error: "", Data: sug})
}
