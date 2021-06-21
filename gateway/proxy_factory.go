package krakend

import (
	jsonschema "github.com/devopsfaith/krakend-jsonschema"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
)

// NewProxyFactory returns a new ProxyFactory wrapping the injected BackendFactory with the default proxy stack and a metrics collector
func NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory) proxy.Factory {
	proxyFactory := proxy.NewDefaultFactory(backendFactory, logger)
	proxyFactory = proxy.NewShadowFactory(proxyFactory)
	proxyFactory = jsonschema.ProxyFactory(proxyFactory)
	return proxyFactory
}

type proxyFactory struct{}

func (p proxyFactory) NewProxyFactory(logger logging.Logger, backendFactory proxy.BackendFactory) proxy.Factory {
	return NewProxyFactory(logger, backendFactory)
}
