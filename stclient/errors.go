package stclient

import (
	"github.com/go-resty/resty/v2"
)

type apiError struct {
	httpStatusCode int
	err            error
}

func (a *apiError) Error() string {
	return "failed to request syncthing api"
}

func (a *apiError) Unwrap() error {
	return a.err
}

func newApiError(err error) *apiError {
	return &apiError{
		err: err,
	}
}

func newHttpApiError(resp *resty.Response) *apiError {
	return &apiError{
		httpStatusCode: resp.StatusCode(),
	}
}
