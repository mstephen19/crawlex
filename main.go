package main

import (
	"fmt"

	"github.com/mstephen19/crawlex/core"
)

func main() {
	requests := []*core.RequestOptions{{
		Url:   "http://google.com",
		Label: "google",
	}, {
		Url: "http://yahoo.com",
	}}

	router := core.NewRouter()

	router.AddDefaultHandler(func(ctx *core.HandlerContext, err error) {
		fmt.Println(ctx.Response.Status)
	})

	router.AddHandler("google", func(ctx *core.HandlerContext, err error) {
		fmt.Println("Requested Google.")
	})

	crawler := core.NewCrawler(&core.CrawlerConfig{
		Handler: router.Handler(),
	})

	crawler.Run(requests...)
}
