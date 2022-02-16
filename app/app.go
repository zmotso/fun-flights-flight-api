package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zmotso/fun-flights-flight-api/controllers"
	"github.com/zmotso/fun-flights-flight-api/httpclient"
	"github.com/zmotso/fun-flights-flight-api/services/routeservice"
)

// Run is the App Entry Point
func Run() {
	e := echo.New()

	conf, err := NewConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Controllers
	routesCtl := controllers.NewRouteController(
		routeservice.NewRouteService(
			conf.FlyProviders,
			httpclient.NewHTTPClient(),
		),
	)

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Fun flights API")
	})
	e.GET("/routes", routesCtl.Get)

	// Start server
	go func() {
		if err := e.Start(conf.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
