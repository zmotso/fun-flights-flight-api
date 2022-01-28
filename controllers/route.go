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

func NewRouteController(routesService routeservice.RouteService) *routeController {
	return &routeController{
		routesService: routesService,
	}
}

func (ctl *routeController) Get(c echo.Context) error {
	routes, err := ctl.routesService.GetRoutes()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, routes)
}
