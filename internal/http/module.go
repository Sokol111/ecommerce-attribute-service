package http

import (
	"github.com/Sokol111/ecommerce-attribute-service-api/gen/httpapi"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
)

func NewHttpHandlerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			newAttributeHandler,
			newOgenServer,
		),
		fx.Invoke(registerOgenRoutes),
	)
}

func newOgenServer(handler httpapi.Handler, tracerProvider trace.TracerProvider, meterProvider metric.MeterProvider) (*httpapi.Server, error) {
	return httpapi.NewServer(
		handler,
		httpapi.WithTracerProvider(tracerProvider),
		httpapi.WithMeterProvider(meterProvider),
	)
}

func registerOgenRoutes(engine *gin.Engine, server *httpapi.Server) {
	// Mount ogen server for API versioned paths only
	// This avoids conflicts with health routes (/health/ready, /health/live)
	engine.Any("/v1/*path", gin.WrapH(server))
	// Add more versions as needed:
	// engine.Any("/v2/*path", gin.WrapH(server))
}
