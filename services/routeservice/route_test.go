package routeservice

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
)

type httpClientMock struct {
	DoMock func(req *retryablehttp.Request) (*http.Response, error)
}

func (c *httpClientMock) Do(req *retryablehttp.Request) (*http.Response, error) {
	return c.DoMock(req)
}

func TestGetRoutes(t *testing.T) {
	client := &httpClientMock{}
	client.DoMock = func(req *retryablehttp.Request) (*http.Response, error) {

		datapr1, err := json.Marshal([]Route{
			{
				Airline:            "air1",
				SourceAirport:      "LV",
				DestinationAirport: "KY",
				CodeShare:          "code1",
				Stops:              0,
				Equipment:          "eq1",
			},
			{
				Airline:            "air2",
				SourceAirport:      "LV",
				DestinationAirport: "BR",
				CodeShare:          "code2",
				Stops:              1,
				Equipment:          "eq2",
			},
		})

		if err != nil {
			t.Fatalf("failed to convert routes %s", err.Error())
		}

		datapr2, err := json.Marshal([]Route{
			{
				Airline:            "air3",
				SourceAirport:      "LV",
				DestinationAirport: "BR",
				CodeShare:          "code3",
				Stops:              0,
				Equipment:          "eq3",
			},
			{
				Airline:            "air4",
				SourceAirport:      "IY",
				DestinationAirport: "BP",
				CodeShare:          "code4",
				Stops:              0,
				Equipment:          "eq4",
			},
		})

		if err != nil {
			t.Fatalf("failed to convert routes %s", err.Error())
		}

		mockData := map[string][]byte{
			"provider1": datapr1,
			"provider2": datapr2,
		}

		r := io.NopCloser(bytes.NewReader(mockData[req.URL.Host]))
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}

	service := NewRouteService(
		[]string{"http://provider1", "http://provider2"},
		client,
	)

	routes, err := service.GetRoutes(context.Background())
	if assert.NoError(t, err) {
		assert.Equal(t, 3, len(routes))
	}
}
