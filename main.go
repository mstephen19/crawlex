package main

import (
	"io"

	"github.com/mstephen19/crawlex/core"
)

func main() {
	router := core.NewRouter(false)

	router.AddHandler("google", func(ctx *core.HandlerContext, err error) {
		b, _ := io.ReadAll(ctx.Response.Body)
		ctx.Push(string(b))
		ctx.Retry()
	})

	crawler := core.NewCrawler(&core.CrawlerConfig{
		Handler: router.Handler(),
	})

	crawler.Run(&core.RequestOptions{
		Url:   "http://google.com",
		Label: "google",
	})
}
