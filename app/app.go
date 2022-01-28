package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zmotso/fun-flights-flight-api/controllers"
	"github.com/zmotso/fun-flights-flight-api/httpclient"
	"github.com/zmotso/fun-flights-flight-api/services/routeservice"
)

// Run is the App Entry Point
func Run() {
	routesCtl := controllers.NewRouteController(
		routeservice.NewRouteService(
			[]string{
				"http://ase.asmt.live:8000/provider/flights1", //TODO: move to env
				"http://ase.asmt.live:8000/provider/flights2",
			},
			httpclient.NewHTTPClient(),
		),
	)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Fun flights API")
	})
	e.POST("/routes", routesCtl.Get)

	e.Logger.Fatal(e.Start(":8080"))
}
