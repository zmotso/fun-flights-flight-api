package routeservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/zmotso/fun-flights-flight-api/httpclient"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("routeservice")

// Route represents flying route between airports
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
	GetRoutes(ctx context.Context) ([]Route, error)
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
		httpClient:      httpClient,
	}
}

// GetRoutes will return merged routes from providers
// TODO:
// - cache result or store in mongodb
// - filter by source and destination airport
func (s *routeService) GetRoutes(ctx context.Context) ([]Route, error) {
	_, span := tracer.Start(ctx, "GetRoutes")
	defer span.End()

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

	if len(allRoutes) == 1 {
		return allRoutes[0], nil
	}

	return mergeRoutes(allRoutes), nil
}

func (s *routeService) getProviderRoutes(url string) ([]Route, error) {
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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

func mergeRoutes(allRoutes [][]Route) []Route {
	if len(allRoutes) == 0 {
		return []Route{}
	}

	if len(allRoutes) == 1 {
		return allRoutes[0]
	}

	routes := map[string][]Route{}
	for i := len(allRoutes) - 1; i >= 0; i-- {
		for j := 0; j < len(allRoutes[i]); j++ {
			currentRoute := allRoutes[i][j]
			if _, ok := routes[currentRoute.SourceAirport]; ok {
				routeExsists := false
				for k := 0; k < len(routes[currentRoute.SourceAirport]); k++ {
					if routes[currentRoute.SourceAirport][k].DestinationAirport == currentRoute.DestinationAirport {
						routes[currentRoute.SourceAirport][k] = currentRoute
						routeExsists = true
						break
					}
				}
				if !routeExsists {
					routes[currentRoute.SourceAirport] = append(routes[currentRoute.SourceAirport], currentRoute)
				}
				continue
			}
			routes[currentRoute.SourceAirport] = []Route{currentRoute}
		}
	}

	mergedroutes := make([]Route, 0, len(allRoutes[1]))
	for _, v := range routes {
		mergedroutes = append(mergedroutes, v...)
	}

	return mergedroutes
}
