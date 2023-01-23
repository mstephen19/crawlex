package core

import (
	"bytes"
	"net/http"
)

func MakeRequest(opts *RequestOptions, client *http.Client) (response *http.Response, err error) {
	request, err := http.NewRequest(opts.Method, opts.Url, bytes.NewReader(opts.Body))
	for key, value := range opts.Headers {
		request.Header.Add(key, value)
	}

	if err != nil {
		return
	}

	response, err = client.Do(request)
	return
}
