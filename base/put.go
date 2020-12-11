package base

import (
	"bytes"
	"context"
	"net/http"
)

type PutHttpRequestInput struct {
	Body             []byte
	ValidStatusCodes []int
	ValidStatusFunc  ValidStatusFunc
	Uri              string
}

func (i PutHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

func (i PutHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

func (c Client) Put(ctx context.Context, input PutHttpRequestInput) (*http.Response, int, error) {
	var status int
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(input.Body))
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.performRequest(ctx, req, input)
	if err != nil {
		return nil, status, err
	}
	return resp, status, nil
}
