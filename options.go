package Malt

import (
	"net/url"
	"os"
	"time"

	"github.com/taluos/Malt/core/registry"
	"github.com/taluos/Malt/pkg/log"
	malitServer "github.com/taluos/Malt/server"
)

type options struct {
	id        string
	name      string
	endpoints []*url.URL
	tags      []string
	version   string

	metadata map[string]string
	signal   []os.Signal

	logger           log.Logger
	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration

	server []malitServer.Server

	// restserver []restserver.Server
	// rpcserver  []rpcserver.Server
}

type Option func(*options)

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithId(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(Name string) Option {
	return func(o *options) {
		o.name = Name
	}
}
func WithTags(tags []string) Option {
	return func(o *options) {
		o.tags = tags
	}
}

func WithVersion(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

func WithSignal(signal []os.Signal) Option {
	return func(o *options) {
		o.signal = signal
	}
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithRegistrarTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.registrarTimeout = timeout
	}
}

func WithStopTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.stopTimeout = timeout
	}
}

func WithLogger(logger log.Logger) Option {
	return func(o *options) {
		o.logger = logger
	}
}
func WithMetadata(metadata map[string]string) Option {
	return func(o *options) {
		o.metadata = metadata
	}
}

func WithServer(server ...malitServer.Server) Option {
	return func(o *options) {
		o.server = server
	}
}

/*
func WithRESTServer(restserver ...restserver.Server) Option {
	return func(o *options) {
		o.restserver = restserver
	}
}

func WithRPCServer(rpcserver ...rpcserver.Server) Option {
	return func(o *options) {
		o.rpcserver = rpcserver
	}
}
*/
