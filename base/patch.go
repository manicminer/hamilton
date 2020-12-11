package base

import (
	"bytes"
	"context"
	"net/http"
)

type PatchHttpRequestInput struct {
	Body             []byte
	ValidStatusCodes []int
	ValidStatusFunc  ValidStatusFunc
	Uri              string
}

func (i PatchHttpRequestInput) GetValidStatusCodes() []int {
	return i.ValidStatusCodes
}

func (i PatchHttpRequestInput) GetValidStatusFunc() ValidStatusFunc {
	return i.ValidStatusFunc
}

func (c Client) Patch(ctx context.Context, input PatchHttpRequestInput) (*http.Response, int, error) {
	var status int
	url := c.buildUri(input.Uri)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(input.Body))
	if err != nil {
		return nil, status, err
	}
	resp, status, err := c.performRequest(ctx, req, input)
	if err != nil {
		return nil, status, err
	}
	return resp, status, nil
}
