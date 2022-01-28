package httpclient

import (
	"net/http"
)

type HTTPClient interface {
	Get(url string) (resp *http.Response, err error)
}

func NewHTTPClient() HTTPClient {
	return &http.Client{}
}
