package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zmotso/fun-flights-flight-api/services/routeservice"
)

// RouteController interface
type RouteController interface {
	Get(c echo.Context) error
}

type routeController struct {
	routesService routeservice.RouteService
}

//NewRouteController will instantiate RouteController
func NewRouteController(routesService routeservice.RouteService) *routeController {
	return &routeController{
		routesService: routesService,
	}
}

//Get will return routes from providers
func (ctl *routeController) Get(c echo.Context) error {
	routes, err := ctl.routesService.GetRoutes()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, routes)
}
