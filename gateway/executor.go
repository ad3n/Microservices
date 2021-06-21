package krakend

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	cmd "github.com/devopsfaith/krakend-cobra"
	cors "github.com/devopsfaith/krakend-cors/gin"
	gelf "github.com/devopsfaith/krakend-gelf"
	gologging "github.com/devopsfaith/krakend-gologging"
	logstash "github.com/devopsfaith/krakend-logstash"
	"github.com/devopsfaith/krakend-usage/client"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/devopsfaith/krakend/proxy"
	krakendrouter "github.com/devopsfaith/krakend/router"
	router "github.com/devopsfaith/krakend/router/gin"
	server "github.com/devopsfaith/krakend/transport/http/server/plugin"
	"github.com/gin-gonic/gin"
	"github.com/go-contrib/uuid"
)

// NewExecutor returns an executor for the cmd package. The executor initalizes the entire gateway by
// registering the components and composing a RouterFactory wrapping all the middlewares.
func NewExecutor(ctx context.Context) cmd.Executor {
	eb := new(ExecutorBuilder)
	return eb.NewCmdExecutor(ctx)
}

// PluginLoader defines the interface for the collaborator responsible of starting the plugin loaders
type PluginLoader interface {
	Load(folder, pattern string, logger logging.Logger)
}

// EngineFactory returns a gin engine, ready to be passed to the KrakenD RouterFactory
type EngineFactory interface {
	NewEngine(config.ServiceConfig, logging.Logger, io.Writer) *gin.Engine
}

// ProxyFactory returns a KrakenD proxy factory, ready to be passed to the KrakenD RouterFactory
type ProxyFactory interface {
	NewProxyFactory(logging.Logger, proxy.BackendFactory) proxy.Factory
}

// BackendFactory returns a KrakenD backend factory, ready to be passed to the KrakenD proxy factory
type BackendFactory interface {
	NewBackendFactory(context.Context, logging.Logger) proxy.BackendFactory
}

// HandlerFactory returns a KrakenD router handler factory, ready to be passed to the KrakenD RouterFactory
type HandlerFactory interface {
	NewHandlerFactory(logging.Logger) router.HandlerFactory
}

// LoggerFactory returns a KrakenD Logger factory, ready to be passed to the KrakenD RouterFactory
type LoggerFactory interface {
	NewLogger(config.ServiceConfig) (logging.Logger, io.Writer, error)
}

// RunServer defines the interface of a function used by the KrakenD router to start the service
type RunServer func(context.Context, config.ServiceConfig, http.Handler) error

// RunServerFactory returns a RunServer with several wraps around the injected one
type RunServerFactory interface {
	NewRunServer(logging.Logger, router.RunServerFunc) RunServer
}

// ExecutorBuilder is a composable builder. Every injected property is used by the NewCmdExecutor method.
type ExecutorBuilder struct {
	LoggerFactory    LoggerFactory
	PluginLoader     PluginLoader
	EngineFactory    EngineFactory
	ProxyFactory     ProxyFactory
	BackendFactory   BackendFactory
	HandlerFactory   HandlerFactory
	RunServerFactory RunServerFactory

	Middlewares []gin.HandlerFunc
}

type LoggerBuilder struct{}

// NewLogger sets up the logging components as defined at the configuration.
func (LoggerBuilder) NewLogger(cfg config.ServiceConfig) (logging.Logger, io.Writer, error) {
	var writers []io.Writer
	gelfWriter, gelfErr := gelf.NewWriter(cfg.ExtraConfig)
	if gelfErr == nil {
		writers = append(writers, gelfWriterWrapper{gelfWriter})
		gologging.SetFormatterSelector(func(w io.Writer) string {
			switch w.(type) {
			case gelfWriterWrapper:
				return "%{message}"
			default:
				return gologging.DefaultPattern
			}
		})
	}

	logger, gologgingErr := logstash.NewLogger(cfg.ExtraConfig)

	if gologgingErr != nil {
		logger, gologgingErr = gologging.NewLogger(cfg.ExtraConfig, writers...)

		if gologgingErr != nil {
			var err error
			logger, err = logging.NewLogger("DEBUG", os.Stdout, "")
			if err != nil {
				return logger, gelfWriter, err
			}
			logger.Error("unable to create the gologging logger:", gologgingErr.Error())
		}
	}

	return logger, gelfWriter, nil
}

// NewCmdExecutor returns an executor for the cmd package. The executor initalizes the entire gateway by
// delegating most of the tasks to the injected collaborators. They register the components and
// compose a RouterFactory wrapping all the middlewares.
// Every nil collaborator is replaced by the default one offered by this package.
func (e *ExecutorBuilder) NewCmdExecutor(ctx context.Context) cmd.Executor {
	e.checkCollaborators()

	return func(cfg config.ServiceConfig) {
		logger, gelfWriter, gelfErr := e.LoggerFactory.NewLogger(cfg)
		if gelfErr != nil {
			return
		}

		logger.Info("Listening on port:", cfg.Port)

		startReporter(ctx, logger, cfg)

		if cfg.Plugin != nil {
			e.PluginLoader.Load(cfg.Plugin.Folder, cfg.Plugin.Pattern, logger)
		}

		// setup the krakend router
		routerFactory := router.NewFactory(router.Config{
			Engine: e.EngineFactory.NewEngine(cfg, logger, gelfWriter),
			ProxyFactory: e.ProxyFactory.NewProxyFactory(
				logger,
				e.BackendFactory.NewBackendFactory(ctx, logger),
			),
			Middlewares:    e.Middlewares,
			Logger:         logger,
			HandlerFactory: e.HandlerFactory.NewHandlerFactory(logger),
			RunServer:      router.RunServerFunc(e.RunServerFactory.NewRunServer(logger, krakendrouter.RunServer)),
		})

		// start the engines
		routerFactory.NewWithContext(ctx).Run(cfg)
	}
}

func (e *ExecutorBuilder) checkCollaborators() {
	if e.PluginLoader == nil {
		e.PluginLoader = new(pluginLoader)
	}

	if e.EngineFactory == nil {
		e.EngineFactory = new(engineFactory)
	}
	if e.ProxyFactory == nil {
		e.ProxyFactory = new(proxyFactory)
	}
	if e.BackendFactory == nil {
		e.BackendFactory = new(backendFactory)
	}
	if e.HandlerFactory == nil {
		e.HandlerFactory = new(handlerFactory)
	}
	if e.LoggerFactory == nil {
		e.LoggerFactory = new(LoggerBuilder)
	}
	if e.RunServerFactory == nil {
		e.RunServerFactory = new(DefaultRunServerFactory)
	}
}

// DefaultRunServerFactory creates the default RunServer by wrapping the injected RunServer
// with the plugin loader and the CORS module
type DefaultRunServerFactory struct{}

func (d *DefaultRunServerFactory) NewRunServer(l logging.Logger, next router.RunServerFunc) RunServer {
	return RunServer(server.New(
		l,
		server.RunServer(cors.NewRunServer(cors.RunServer(next))),
	))
}

const (
	usageDisable = "USAGE_DISABLE"
)

func startReporter(ctx context.Context, logger logging.Logger, cfg config.ServiceConfig) {
	if os.Getenv(usageDisable) == "1" {
		logger.Info("usage report client disabled")
		return
	}

	clusterID, err := cfg.Hash()
	if err != nil {
		logger.Warning("unable to hash the service configuration:", err.Error())
		return
	}

	go func() {
		serverID := uuid.NewV4().String()
		logger.Info(fmt.Sprintf("registering usage stats for cluster ID '%s'", clusterID))

		if err := client.StartReporter(ctx, client.Options{
			ClusterID: clusterID,
			ServerID:  serverID,
		}); err != nil {
			logger.Warning("unable to create the usage report client:", err.Error())
		}
	}()
}

type gelfWriterWrapper struct {
	io.Writer
}
