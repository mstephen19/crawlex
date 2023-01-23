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
	Crawler  *Crawler
}

func (ctx *HandlerContext) Enqueue(requests ...*RequestOptions) error {
	return ctx.Crawler.Enqueue(requests...)
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
