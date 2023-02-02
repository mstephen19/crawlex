package core

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type HandlerContext struct {
	Options  *RequestOptions
	Response *http.Response
	crawler  *Crawler
	proxy    *Proxy
}

func (ctx *HandlerContext) Proxy() string {
	if ctx.proxy == nil {
		return ""
	}

	return ctx.proxy.raw
}

func (ctx *HandlerContext) Push(data ...any) {
	ctx.crawler.store.Push(data...)
}

func (ctx *HandlerContext) MarkProxyBad() {
	ctx.crawler.proxyPool.MarkBad(ctx.proxy)
}

func (ctx *HandlerContext) MarkProxyGood() {
	ctx.crawler.proxyPool.MarkGood(ctx.proxy)
}

func (ctx *HandlerContext) Retry() error {
	return ctx.crawler.Enqueue(ctx.Options)
}

func (ctx *HandlerContext) Enqueue(requests ...*RequestOptions) error {
	return ctx.crawler.Enqueue(requests...)
}

func (ctx *HandlerContext) ParseHTML() (doc *goquery.Document, err error) {
	return goquery.NewDocumentFromReader(ctx.Response.Body)
}

func (ctx *HandlerContext) JSON(target any) (err error) {
	bytes, err := io.ReadAll(ctx.Response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, target)
	return
}

type HandlerFunc func(ctx *HandlerContext, err error)
