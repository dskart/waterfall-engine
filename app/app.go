package app

import (
	"context"
	"fmt"

	"github.com/dskart/waterfall-engine/app/engine"
	"github.com/dskart/waterfall-engine/store"
	"go.uber.org/zap"
)

type App struct {
	store  *store.Store
	config Config
	logger *zap.Logger
	engine engine.Engine
}

type Options struct {
	store *store.Store
}

func WithStore(store *store.Store) func(*Options) {
	return func(o *Options) {
		o.store = store
	}
}

func New(ctx context.Context, logger *zap.Logger, config Config, opts ...func(*Options)) (*App, error) {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}

	if options.store == nil {
		var err error
		options.store, err = store.New(ctx, config.Store)
		if err != nil {
			return nil, fmt.Errorf("could not create store: %w", err)
		}
	}

	engine := engine.NewEngine(config.Engine)

	ret := &App{
		store:  options.store,
		config: config,
		logger: logger,
		engine: engine,
	}

	return ret, nil
}

func (a *App) Close(ctx context.Context) error {
	if a.store != nil {
		if err := a.store.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) Logger() *zap.Logger {
	return a.logger
}
