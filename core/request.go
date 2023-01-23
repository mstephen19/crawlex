package core

import (
	"errors"
	"net/http"
)

type RequestOptions struct {
	Url         string
	Method      string
	Body        []byte
	Headers     map[string]string
	Label       any
	SkipRequest bool
	userData    map[any]any
}

func (options *RequestOptions) Set(key any, value any) {
	options.userData[key] = value
}

func (options *RequestOptions) Get(key any) (value any, exists bool) {
	value, exists = options.userData[key]

	return
}

func CleanRequestOptions(opts *RequestOptions) (err error) {
	if opts.Url == "" {
		err = errors.New("Empty or invalid URL provided in RequestOptions.")
		return
	}

	defaultIfZero(&opts.Method, "", http.MethodGet)

	return
}
