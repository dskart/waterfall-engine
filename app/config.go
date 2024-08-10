package app

import (
	"github.com/dskart/waterfall-engine/app/engine"
	"github.com/dskart/waterfall-engine/store"
)

type Config struct {
	Store  store.Config  `yaml:"Store"`
	Engine engine.Config `yaml:"Engine"`
}
