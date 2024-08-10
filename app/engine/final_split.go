package engine

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

var ErrFinalSplitStageHasRemainingCapital = fmt.Errorf("final split stage has remaining capital")
var ErrFinalSplitStageDistributionNotEqualStartingCapital = fmt.Errorf("final split stage distribution does not equal starting capital")

func calculateFinalSplit(cfg FinalSplitConfig, startingCapital *money.Money) (TierStage, error) {
	allocations, err := startingCapital.Allocate(int(cfg.LpPercentage*100), int(cfg.GpPercentage*100))
	if err != nil {
		return TierStage{}, fmt.Errorf("could not allocation final split")
	} else if len(allocations) != 2 {
		return TierStage{}, fmt.Errorf("len of allocations is not 2")
	}

	lpAllocation := allocations[0]
	gpAllocation := allocations[1]

	distribution, err := lpAllocation.Add(gpAllocation)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not add lpAllocation and gpAllocation: %w", err)
	}
	remainingCapital, err := startingCapital.Subtract(distribution)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not subtract startingCapital and distribution: %w", err)
	}

	ret := TierStage{
		TierName:          FinalSplitStage,
		StartingCapital:   startingCapital,
		LpAllocattion:     lpAllocation,
		GpAllocattion:     gpAllocation,
		TotalDistribution: distribution,
		RemainingCapital:  remainingCapital,
	}

	if !ret.RemainingCapital.IsZero() {
		return TierStage{}, ErrFinalSplitStageHasRemainingCapital
	}
	if ok, err := ret.TotalDistribution.Equals(startingCapital); err != nil {
		return TierStage{}, fmt.Errorf("could not compare ret.TotalDistribution and startingCapital")
	} else if !ok {
		return TierStage{}, ErrFinalSplitStageDistributionNotEqualStartingCapital
	}

	return ret, nil
}
