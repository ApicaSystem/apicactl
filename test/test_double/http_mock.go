package test_double

import (
	"bytes"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HttpMock struct {
	mock.Mock
}

func (mockClient *HttpMock) MakeApiCall(method string, url string, payload *bytes.Buffer) (*http.Response, error) {
	args := mockClient.Called(method, url, payload)

	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).(*http.Response), nil
}
