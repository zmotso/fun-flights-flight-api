package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zmotso/fun-flights-flight-api/controllers"
	"github.com/zmotso/fun-flights-flight-api/httpclient"
	"github.com/zmotso/fun-flights-flight-api/services/routeservice"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Run is the App Entry Point
func Run() {
	e := echo.New()

	conf, err := NewConfig()
	if err != nil {
		e.Logger.Fatal(err)
	}

	tp := initTracer(conf)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Controllers
	routesCtl := controllers.NewRouteController(
		routeservice.NewRouteService(
			conf.FlyProviders,
			httpclient.NewHTTPClient(e.Logger),
		),
	)

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware("flight-server"))

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Fun flights API")
	})
	e.GET("/routes", routesCtl.Get)

	// Start server
	go func() {
		if err := e.Start(conf.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server %s", err.Error())
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

func initTracer(conf *Config) *sdktrace.TracerProvider {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(conf.JaegerURL)))
	if err != nil {
		log.Fatal(err)
	}

	resource, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fun-flights-flight-api"),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp
}
