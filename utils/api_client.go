package utils

import (
	"bytes"
	"fmt"
	"net/http"
)

type Api interface {
	MakeApiCall(method string, url string, payload *bytes.Buffer) (*http.Response, error)
}

type ApiClient struct{}

func (c *ApiClient) MakeApiCall(method string, url string, payload *bytes.Buffer) (*http.Response, error) {
	var req *http.Request
	var err error
	var res *http.Response

	if method == http.MethodPost {
		req, err = CreateHttpRequest(method, url, payload)
	} else {
		req, err = CreateHttpRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("Error sending request")
	}

	client := http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err.Error())
	}

	return res, nil
}
