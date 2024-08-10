package engine

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

func calculateReturnOfCapital(startingCapital *money.Money, contributions []Contribution) (TierStage, []Contribution, error) {
	capitalLeft := startingCapital
	for i, c := range contributions {
		if capitalLeft.IsZero() {
			break
		} else if c.ReturnCapitalLeft.IsZero() {
			continue
		}

		if ok, err := c.ReturnCapitalLeft.GreaterThanOrEqual(capitalLeft); err != nil {
			return TierStage{}, nil, fmt.Errorf("could not compare roc and capital left: %w", err)
		} else if ok {
			c.ReturnCapitalLeft, err = c.ReturnCapitalLeft.Subtract(capitalLeft)
			if err != nil {
				return TierStage{}, nil, fmt.Errorf("could not substract c.ReturnCapitalLeft and capitalLeft: %w", err)
			}
			capitalLeft = money.New(0, money.USD)
		} else {
			capitalLeft, err = capitalLeft.Subtract(c.ReturnCapitalLeft)
			if err != nil {
				return TierStage{}, nil, fmt.Errorf("could not substract capitalLeft and c.RetrunCapitalLeft: %w", err)
			}
			c.ReturnCapitalLeft = money.New(0, money.USD)
		}
		contributions[i] = c
	}

	distribution, err := startingCapital.Subtract(capitalLeft)
	if err != nil {
		return TierStage{}, nil, fmt.Errorf("could not subtract startingCapital and capitalLeft: %w", err)
	}
	remainingCapital, err := startingCapital.Subtract(distribution)
	if err != nil {
		return TierStage{}, nil, fmt.Errorf("could not subtract startingCapital and distribution: %w", err)
	}
	ret := TierStage{
		TierName:          RocTierStage,
		StartingCapital:   startingCapital,
		LpAllocattion:     distribution,
		GpAllocattion:     money.New(0, money.USD),
		TotalDistribution: distribution,
		RemainingCapital:  remainingCapital,
	}

	return ret, contributions, nil
}
