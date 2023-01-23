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
