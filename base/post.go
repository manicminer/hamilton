package base

import (
	"bytes"
	"context"
	"net/http"
)

type PostHttpRequestInput struct {
	Body             []byte
	ValidStatusCodes []int
	ValidStatusFunc  ValidStatusFunc
	Uri              string
}

func (i PostHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

func (i PostHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

func (c Client) Post(ctx context.Context, input PostHttpRequestInput) (*http.Response, int, error) {
	var status int
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(input.Body))
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.performRequest(ctx, req, input)
	if err != nil {
		return nil, status, err
	}
	return resp, status, nil
}
