package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis-force/tisearch/server/model"
)

func (h *Handler) SearchUser(c *gin.Context) {
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		c.AbortWithStatusJSON(400, model.UserSearchResponse{Code: 400, Error: "keyword 不能为空"})
		return
	}
	data, err := h.srv.SearchUser(keyword)
	if err != nil {
		c.AbortWithStatusJSON(500, model.UserSearchResponse{Code: 500, Error: err.Error()})
		return
	}
	c.JSON(200, model.UserSearchResponse{Code: 0, Error: "", Data: data})
}

func (h *Handler) SuggestUser(c *gin.Context) {
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		c.AbortWithStatusJSON(400, model.UserSuggestionResponse{Code: 400, Error: "keyword 不能为空"})
		return
	}
	data, err := h.srv.SuggestUser(keyword)
	if err != nil {
		c.AbortWithStatusJSON(500, model.UserSuggestionResponse{Code: 500, Error: err.Error()})
		return
	}
	sug := model.Suggestion{Suggestion: data}
	c.JSON(200, model.TweetSuggestionResponse{Code: 0, Error: "", Data: sug})
}
