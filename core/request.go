package core

import (
	"errors"
	"net/http"
)

type RequestOptions struct {
	Url     string
	Method  string
	Body    []byte
	Headers map[string]string
	Label   any
}

func CleanRequestOptions(opts *RequestOptions) (err error) {
	if opts.Url == "" {
		err = errors.New("Empty or invalid URL provided in RequestOptions.")
		return
	}

	defaultIfZero(&opts.Method, "", http.MethodGet)

	return
}
