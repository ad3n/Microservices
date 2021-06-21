package krakend

import (
	botdetector "github.com/devopsfaith/krakend-botdetector/gin"
	juju "github.com/devopsfaith/krakend-ratelimit/juju/router/gin"
	"github.com/devopsfaith/krakend/logging"
	router "github.com/devopsfaith/krakend/router/gin"
)

// NewHandlerFactory returns a HandlerFactory with a rate-limit and a metrics collector middleware injected
func NewHandlerFactory(logger logging.Logger) router.HandlerFactory {
	handlerFactory := juju.HandlerFactory
	handlerFactory = botdetector.New(handlerFactory, logger)
	return handlerFactory
}

type handlerFactory struct{}

func (h handlerFactory) NewHandlerFactory(l logging.Logger) router.HandlerFactory {
	return NewHandlerFactory(l)
}
