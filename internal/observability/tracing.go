package observability

import (
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func InitTracing(logger *zap.Logger, serviceName string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)
	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}
