package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/redis-force/tisearch/server/service"
)

type Handler struct {
	srv *service.TiSearchService
}

func NewHandler(srv *service.TiSearchService) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) SearchTweet(c *gin.Context) {
}
func (h *Handler) SuggestTweet(c *gin.Context) {
}
