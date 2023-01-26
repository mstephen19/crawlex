package core

import (
	"bytes"
	"io"
	"net/http"
)

func MakeRequest(opts *RequestOptions, client *http.Client) (response *http.Response, err error) {
	body := func() io.Reader {
		if len(opts.Body) > 0 {
			return bytes.NewReader(opts.Body)
		}
		return nil
	}()

	request, err := http.NewRequest(opts.Method, opts.Url, body)
	for key, value := range opts.Headers {
		request.Header.Add(key, value)
	}

	if err != nil {
		return
	}

	response, err = client.Do(request)
	return
}
