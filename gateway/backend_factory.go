package krakend

import (
	"context"

	httpcache "github.com/devopsfaith/krakend-httpcache"
	martian "github.com/devopsfaith/krakend-martian"
	opencensus "github.com/devopsfaith/krakend-opencensus"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/proxy"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	"github.com/devopsfaith/krakend/transport/http/client"
	httprequestexecutor "github.com/devopsfaith/krakend/transport/http/client/plugin"
)

func NewBackendFactory(logger logging.Logger) proxy.BackendFactory {
	return NewBackendFactoryWithContext(context.Background(), logger)
}

func NewBackendFactoryWithContext(ctx context.Context, logger logging.Logger) proxy.BackendFactory {
	requestExecutorFactory := func(cfg *config.Backend) client.HTTPRequestExecutor {
		return opencensus.HTTPRequestExecutor(httpcache.NewHTTPClient(cfg))
	}

	requestExecutorFactory = httprequestexecutor.HTTPRequestExecutor(logger, requestExecutorFactory)
	backendFactory := martian.NewConfiguredBackendFactory(logger, requestExecutorFactory)
	backendFactory = juju.BackendFactory(backendFactory)
	backendFactory = opencensus.BackendFactory(backendFactory)

	return backendFactory
}

type backendFactory struct{}

func (b backendFactory) NewBackendFactory(ctx context.Context, l logging.Logger) proxy.BackendFactory {
	return NewBackendFactoryWithContext(ctx, l)
}
