package cmd

import (
	"github.com/dskart/waterfall-engine/app"
)

const cfgEnvPrefix = "BD"

type Config struct {
	App app.Config `yaml:"App"`
}
