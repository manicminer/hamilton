package base

import (
	"context"
	"net/http"
)

type DeleteHttpRequestInput struct {
	ValidStatusCodes []int
	ValidStatusFunc  ValidStatusFunc
	Uri              string
}

func (i DeleteHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

func (i DeleteHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

func (c Client) Delete(ctx context.Context, input DeleteHttpRequestInput) (*http.Response, int, error) {
	var status int
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, http.NoBody)
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.performRequest(ctx, req, input)
	if err != nil {
		return nil, status, err
	}
	return resp, status, nil
}
