package engine

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

func calculateCatchUp(cfg CatchUpConfig, startingCapital *money.Money, preferredReturn *money.Money) (TierStage, error) {
	cu := cfg.CatchupPercentage
	ci := cfg.CarriedInterestPercentage
	totalPercentage := ci / (cu - ci)
	rawCatchupAmount := preferredReturn.AsMajorUnits() * totalPercentage
	catchupAmount := money.NewFromFloat(rawCatchupAmount, money.USD)

	capitalLeft := startingCapital
	if ok, err := catchupAmount.GreaterThanOrEqual(capitalLeft); err != nil {
		return TierStage{}, fmt.Errorf("could not compare catchupAmount and startingCapital: %w", err)
	} else if ok {
		capitalLeft = money.New(0, money.USD)
	} else {
		var err error
		capitalLeft, err = capitalLeft.Subtract(catchupAmount)
		if err != nil {
			return TierStage{}, fmt.Errorf("could not subtract capitalLeft and catchupAmount: %w", err)
		}
	}

	totalCatchupDistribution, err := startingCapital.Subtract(capitalLeft)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not subtract startingCapital and capitalLeft: %w", err)
	}
	remainingCapital, err := startingCapital.Subtract(totalCatchupDistribution)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not subtract startingCapital and distribution: %w", err)
	}

	gpDistribution := totalCatchupDistribution * cfg.CatchupPercentage
	lpDistribution := totalCatchupDistribution * (1 - cfg.CatchupPercentage)

	ret := TierStage{
		TierName:          CatchUpStage,
		StartingCapital:   startingCapital,
		LpAllocattion:     money.New(0, money.USD),
		GpAllocattion:     gpDistribution,
		TotalDistribution: lpDistribution,
		RemainingCapital:  remainingCapital,
	}

	return ret, nil
}
