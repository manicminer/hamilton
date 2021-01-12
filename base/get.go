package base

import (
	"context"
	"fmt"
	"net/http"
)

type GetHttpRequestInput struct {
	ValidStatusCodes []int
	ValidStatusFunc  ValidStatusFunc
	Uri              Uri
}

func (i GetHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

func (i GetHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

func (c Client) Get(ctx context.Context, input GetHttpRequestInput) (*http.Response, int, error) {
	var status int
	url, err := c.buildUri(input.Uri)
	if err != nil {
		return nil, status, fmt.Errorf("unable to make request: %v", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.performRequest(req, input)
	if err != nil {
		return nil, status, err
	}
	return resp, status, nil
}
