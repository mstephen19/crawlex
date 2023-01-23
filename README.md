# Crawlex

This is a passion project inspired by the [**Crawlee** TypeScript library](https://crawlee.dev/). It will potentially be published in the future, but right now it is just a proof-of-concept.

## Usage

Requests are interpreted by Crawlex via `RequestOptions`.

```go
requests := []*core.RequestOptions{{
    Url:   "http://google.com",
    Label: "google",
}, {
    Url: "http://yahoo.com",
}}
```

To handle requests, you can either create your own `HandlerFunc`, or create a router that handles different paths via labelled requests.

```go
router := core.NewRouter()

router.AddDefaultHandler(func(ctx *core.HandlerContext, err error) {
    fmt.Println(ctx.Response.Status)
})

router.AddHandler("google", func(ctx *core.HandlerContext, err error) {
    fmt.Println("Requested Google.")
})
```

A crawler can then be created. The initial requests can be passed into `Run`.

```go
crawler := core.NewCrawler(&core.CrawlerConfig{
    Handler: router.Handler(),
})

crawler.Run(requests...)
```

Final code:

```go
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
```
