package engine

import (
	"fmt"
	"math"
	"time"

	"github.com/Rhymond/go-money"
)

var ErrRocLeft = fmt.Errorf("contribution still has roc left")

func calculatePreferredReturn(cfg PreferredReturnConfig, date time.Time, startingCapital *money.Money, contributions []Contribution) (TierStage, error) {
	capitalLeft := startingCapital
	for _, c := range contributions {
		if capitalLeft.IsZero() {
			break
		} else if !c.ReturnCapitalLeft.IsZero() {
			return TierStage{}, ErrRocLeft
		}
		cashFlow := c.Amount.AsMajorUnits()
		timeDiff := date.Sub(c.Date)
		days := int64(timeDiff.Hours() / 24)
		exponent := float64(days) / 365.0
		rawPreferredReturn := cashFlow*math.Pow((1+cfg.HurdlePercentage), exponent) - cashFlow
		preferredReturn := money.NewFromFloat(rawPreferredReturn, money.USD)

		if ok, err := preferredReturn.GreaterThanOrEqual(capitalLeft); err != nil {
			return TierStage{}, fmt.Errorf("could not compare preferredReturn and capitalleft: %w", err)
		} else if ok {
			capitalLeft = money.New(0, money.USD)
		} else {
			var err error
			capitalLeft, err = capitalLeft.Subtract(preferredReturn)
			if err != nil {
				return TierStage{}, fmt.Errorf("could not subtract capitalLeft and preferedReturn: %w", err)
			}
		}
	}

	distribution, err := startingCapital.Subtract(capitalLeft)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not subtract startingCapital and capitalLeft: %w", err)
	}
	remainingCapital, err := startingCapital.Subtract(distribution)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not subtract startingCapital and distribution: %w", err)
	}
	ret := TierStage{
		TierName:          PreferredReturnStage,
		StartingCapital:   startingCapital,
		LpAllocattion:     distribution,
		GpAllocattion:     money.New(0, money.USD),
		TotalDistribution: distribution,
		RemainingCapital:  remainingCapital,
	}
	return ret, nil
}
