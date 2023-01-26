package main

import (
	"fmt"

	"github.com/mstephen19/crawlex/core"
)

func main() {
	requests := []*core.RequestOptions{{
		Url:   "http://google.com",
		Label: "google",
	}}

	router := core.NewRouter(false)

	router.AddHandler("google", func(ctx *core.HandlerContext, err error) {
		fmt.Println("First handler.")
	})

	router.AddHandler("google", func(ctx *core.HandlerContext, err error) {
		fmt.Println("Second handler.")
	})

	crawler := core.NewCrawler(&core.CrawlerConfig{
		Handler: router.Handler(),
	})

	crawler.Run(requests...)
}
