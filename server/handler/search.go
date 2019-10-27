package handler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis-force/tisearch/server/model"
	"github.com/redis-force/tisearch/server/service"
)

var (
	regex = regexp.MustCompile(`.*AGAINST\s{0,}\('(.*)'.*\).*`)
)

func mathQuery(q string) (table, keywrod, query string, err error) {
	if !strings.HasPrefix(q, "SELECT") {
		query = q
		return
	}
	qs := strings.SplitN(q, "WHERE", 2)
	selects := strings.TrimSpace(strings.TrimPrefix(qs[0], "SELECT"))
	wheres := strings.SplitN(qs[1], "AND", 2)
	ands := ""
	if len(wheres) == 2 {
		ands = wheres[1]
	}
	if strings.Index(q, "users") != -1 {
		table = "users"
	} else {
		table = "tweets"
	}
	query = selects
	if ands != "" {
		query += " WHERE" + ands
	}
	found := regex.FindStringSubmatch(q)
	fmt.Println(found)
	if len(found) == 2 {
		keywrod = found[1]
		keywrod = strings.Replace(strings.Replace(keywrod, " and ", " AND ", -1), ":", "_", -1)
	}
	return
}

type Handler struct {
	srv *service.TiSearchService
}

func NewHandler(srv *service.TiSearchService) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) Search(c *gin.Context) {
	q := c.Query("q")
	if len(q) == 0 {
		c.AbortWithStatusJSON(400, model.SearchResponse{Code: 400, Error: "q不能为空"})
		return
	}
	tp, keyword, query, err := mathQuery(q)
	if err != nil {
		c.AbortWithStatusJSON(400, model.SearchResponse{Code: 400, Error: err.Error()})
		return
	}
	var data interface{}
	var plans []model.SQLPlan
	rowAffected := int64(-1)
	if tp == "tweets" {
		data, plans, err = h.srv.SearchTweet(keyword, query)
	} else if tp == "users" {
		data, plans, err = h.srv.SearchUser(keyword, query)
	} else {
		data = make([]interface{}, 0)
		rowAffected, plans, err = h.srv.Execute(query)
	}
	if err != nil {
		c.AbortWithStatusJSON(500, model.SearchResponse{Code: 500, Error: err.Error()})
		return
	}
	if rowAffected != -1 {
		c.JSON(200, model.SearchResponse{Code: 204, Error: "", Data: data, Type: tp, Plans: plans, RowAffected: rowAffected})
		return
	}
	c.JSON(200, model.SearchResponse{Code: 0, Error: "", Data: data, Type: tp, Plans: plans, RowAffected: rowAffected})
}
