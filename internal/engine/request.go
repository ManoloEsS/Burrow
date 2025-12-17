package engine

import (
	"context"
	"io"
	"net/http"
)

func GetRequest(ctx *context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequestWithContext(context.Background(), method, url, body)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
