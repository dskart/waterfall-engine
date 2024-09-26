package engine

type Config struct {
	PreferredReturn PreferredReturnConfig `yaml:"PreferredReturn"`
	CatchUp         CatchUpConfig         `yaml:"CatchUp"`
	FinalSplit      FinalSplitConfig      `yaml:"FinalSplit"`
}

type PreferredReturnConfig struct {
	HurdlePercentage float64 `yaml:"HurdlePercentage"`
}

type CatchUpConfig struct {
	Enabled                   bool
	CatchupPercentage         float64 `yaml:"CatchupPercentage"`
	CarriedInterestPercentage float64 `yaml:"CarriedInterestPercentage"`
}

type FinalSplitConfig struct {
	LpPercentage float64 `yaml:"LpPercentage"`
	GpPercentage float64 `yaml:"GpPercentage"`
}
