package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/redis-force/tisearch/server/handler"
	"github.com/redis-force/tisearch/server/service"
)

func main() {
	srv, err := service.NewSearchService()
	if err != nil {
		panic(err)
	}
	hdl := handler.NewHandler(srv)

	mux := gin.New()
	mux.Use(gin.Logger(), gin.Recovery())
	mux.GET("/api/v1/search", hdl.Search)
	// {

	// search.GET("/user", hdl.SearchUser)
	// search.GET("/tweet", hdl.SearchTweet)
	// }
	suggest := mux.Group("/api/v1/suggest")
	{
		suggest.GET("/user", hdl.SuggestUser)
		suggest.GET("/tweet", hdl.SuggestTweet)
	}

	fmt.Println(mux.Run("0.0.0.0:8080"))
}
