package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	defaultTimeOut = time.Second * 10
)

type HTTPClient interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
}

type errResp struct {
	Error string `json:"error"`
}

func NewHTTPClient(logger interface{}) HTTPClient {
	client := retryablehttp.NewClient()
	client.HTTPClient.Timeout = defaultTimeOut
	client.Logger = logger
	client.CheckRetry = checkRetry
	client.HTTPClient.Transport = otelhttp.NewTransport(client.HTTPClient.Transport)

	return client
}

func checkRetry(ctx context.Context, resp *http.Response, err error) (bool, error) {
	retry, err := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

	if !retry {
		buf := bytes.NewBuffer(make([]byte, 0))

		tee := io.TeeReader(resp.Body, buf)

		var reqBody errResp
		_ = json.NewDecoder(tee).Decode(&reqBody)

		defer resp.Body.Close()

		resp.Body = ioutil.NopCloser(buf)

		if reqBody.Error != "" {
			return true, fmt.Errorf("response contains error %s", reqBody.Error)
		}
	}

	return retry, err
}
