package main

import (
	"fmt"

	"github.com/mstephen19/crawlex/core"
)

func main() {
	router := core.NewRouter()

	router.AddDefaultHandler(func(ctx *core.HandlerContext, err error) {
		fmt.Println(ctx.Options.Url, ctx.Proxy())
		ctx.MarkProxyBad()
		ctx.Retry()
	})

	crawler := core.NewCrawler(&core.CrawlerConfig{
		Handler:               router.Handler(),
		Proxies:               []string{"http://125.141.139.198:5566"},
		RequestTimeoutSeconds: 2,
	})

	crawler.Run(&core.RequestOptions{
		Url: "http://typescriptlang.org/",
	})
}
