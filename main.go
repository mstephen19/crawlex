package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/mstephen19/crawlex/core"
)

func main() {
	opts := &core.RequestOptions{
		Url: "http://crawlee.dev/",
	}

	config := &core.CrawlerConfig{
		MaxConcurrency: 50,
		Handler: func(ctx *core.HandlerContext, err error) {
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			doc, err := ctx.ParseHTML()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println(ctx.Options.Url)

			hrefs := doc.Find("a[href]")

			for _, node := range hrefs.Nodes {
				href := func() (result string) {
					for _, attr := range node.Attr {
						if attr.Key == "href" {
							result = attr.Val
							return
						}
					}
					return
				}()

				if strings.HasPrefix(href, "http") {
					continue
				}

				joined, _ := url.JoinPath(ctx.Options.Url, href)
				ctx.Enqueue(&core.RequestOptions{
					Url: joined,
				})
			}
		},
	}

	crawler := core.NewCrawler(config)

	crawler.Run(opts, opts, opts)
}
