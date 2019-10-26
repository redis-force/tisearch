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
	search := mux.Group("/api/v1/search")
	{

		search.GET("/album", hdl.SearchAlbum)
		search.GET("/tweet", hdl.SearchTweet)
	}
	suggest := mux.Group("/api/v1/suggest")
	{
		suggest.GET("/album", hdl.SuggestAlbum)
		suggest.GET("/tweet", hdl.SuggestTweet)
	}

	fmt.Println(mux.Run("0.0.0.0:8080"))
}
