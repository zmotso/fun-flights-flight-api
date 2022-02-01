package app

import (
	"net/http"
	"os"
	"strings"

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
			strings.Split(os.Getenv("FLY_PROVIDERS"), ","),
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
	e.GET("/routes", routesCtl.Get)

	e.Logger.Fatal(e.Start(":8080"))
}
