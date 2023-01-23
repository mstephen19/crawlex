package core

import (
	"bytes"
	"net/http"
)

func MakeRequest(opts *RequestOptions) (response *http.Response, err error) {
	request, err := http.NewRequest(opts.Method, opts.Url, bytes.NewReader(opts.Body))
	for key, value := range opts.Headers {
		request.Header.Add(key, value)
	}

	if err != nil {
		return
	}

	response, err = http.DefaultClient.Do(request)
	return
}
