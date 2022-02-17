package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClientGetWithRetries(t *testing.T) {
	type respOK struct {
		Data string `json:"data"`
	}
	expecedRespOK := respOK{Data: "some data"}

	respIsOK := true
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !respIsOK {
			respErr, err := json.Marshal(map[string]string{"error": "some error"})
			if err != nil {
				t.Fatalf("failed to convert response %s", err.Error())
			}
			if _, err := w.Write(respErr); err != nil {
				t.Fatalf("failed to write response %s", err.Error())
			}
			respIsOK = true
			return
		}

		respSuccess, err := json.Marshal(expecedRespOK)
		if err != nil {
			t.Fatalf("failed to convert response %s", err.Error())
		}
		if _, err := w.Write(respSuccess); err != nil {
			t.Fatalf("failed to write response %s", err.Error())
		}
	}))
	defer ts.Close()

	client := NewHTTPClient(nil).(*retryablehttp.Client)
	client.RetryWaitMin = 10 * time.Millisecond
	client.RetryWaitMax = 10 * time.Millisecond
	client.RetryMax = 2

	resp, err := client.Get(ts.URL)
	if assert.NoError(t, err) {
		respOKGot := respOK{}
		err = json.NewDecoder(resp.Body).Decode(&respOKGot)
		if assert.NoError(t, err) {
			assert.Equal(t, expecedRespOK.Data, respOKGot.Data)
		}
	}
}
