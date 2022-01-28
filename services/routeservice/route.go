package routeservice

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/zmotso/fun-flights-flight-api/httpclient"
)

type Route struct {
	Airline            string `json:"airline"`
	SourceAirport      string `json:"sourceAirport"`
	DestinationAirport string `json:"destinationAirport"`
	CodeShare          string `json:"codeShare"`
	Stops              int    `json:"stops"`
	Equipment          string `json:"equipment"`
}

// RouteService interface
type RouteService interface {
	GetRoutes() ([]Route, error)
}

type routeService struct {
	routesProviders []string
	httpClient      httpclient.HTTPClient
}

// NewRouteService will instantiate RouteService
func NewRouteService(
	routesProviders []string,
	httpClient httpclient.HTTPClient,
) RouteService {
	return &routeService{
		routesProviders: routesProviders,
	}
}

func (s *routeService) GetRoutes() ([]Route, error) {
	ch := make(chan []Route)

	for _, providerURL := range s.routesProviders {
		go func(providerURL string) {
			routes, err := s.getProviderRoutes(providerURL)
			if err != nil {
				log.Println("error getting routes", providerURL, err.Error())
			}
			ch <- routes
		}(providerURL)
	}

	var allRoutes [][]Route
	for i := 0; i < len(s.routesProviders); i++ {
		if routes := <-ch; routes != nil {
			allRoutes = append(allRoutes, routes)
		}
	}

	return nil, nil
}

func (s *routeService) getProviderRoutes(url string) ([]Route, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("failed to get routes, body: %s", string(body))
	}

	routes := []Route{}

	err = json.NewDecoder(resp.Body).Decode(&routes)
	if err != nil {
		return nil, err
	}

	return routes, nil
}
