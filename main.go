package main

import (
	"fmt"

	"github.com/mstephen19/crawlex/core"
)

func main() {
	router := core.NewRouter()

	router.AddDefaultHandler(func(ctx *core.HandlerContext, err error) {
		fmt.Println("Unknown route reached")
	})

	crawler := core.NewCrawler(&core.CrawlerConfig{
		MaxConcurrency: 100,
		Handler:        router.Handler(),
	})

	crawler.Run(&core.RequestOptions{
		Url: "http://www.typescriptlang.org/",
	})
}
