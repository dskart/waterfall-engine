package engine

type Config struct {
	PreferedReturn PreferedReturnConfig
	CatchUp        CatchUpConfig
	FinalSplit     FinalSplitConfig
}

type PreferedReturnConfig struct {
	HurdlePercentage float64
}

type CatchUpConfig struct {
	CatchupPercentage        float64
	CariedInterestPercentage float64
}

type FinalSplitConfig struct {
	LpPercentage float64
	GpPercentage float64
}
