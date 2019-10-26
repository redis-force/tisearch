package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis-force/tisearch/server/model"
)

func (h *Handler) SearchUser(c *gin.Context) {
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		c.AbortWithStatusJSON(400, model.SearchResponse{Code: 400, Error: "keyword 不能为空"})
		return
	}
	data, _, err := h.srv.SearchUser(keyword)
	if err != nil {
		c.AbortWithStatusJSON(500, model.SearchResponse{Code: 500, Error: err.Error()})
		return
	}
	c.JSON(200, model.SearchResponse{Code: 0, Error: "", Data: data})
}

func (h *Handler) SuggestUser(c *gin.Context) {
	keyword := c.Query("keyword")
	if len(keyword) == 0 {
		c.AbortWithStatusJSON(400, model.SuggestionResponse{Code: 400, Error: "keyword 不能为空"})
		return
	}
	data, err := h.srv.SuggestUser(keyword)
	if err != nil {
		c.AbortWithStatusJSON(500, model.SuggestionResponse{Code: 500, Error: err.Error()})
		return
	}
	sug := model.Suggestion{Suggestion: data}
	c.JSON(200, model.SuggestionResponse{Code: 0, Error: "", Data: sug})
}
