package app

import "github.com/dskart/waterfall-engine/store"

type Config struct {
	Store store.Config `yaml:"Store"`
}

func (c Config) Validate() error {
	return nil
}
