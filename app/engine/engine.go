package engine

import (
	"fmt"
	"math"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/dskart/waterfall-engine/model"
)

type Engine struct {
}

type State struct {
	LpContributions []Contribution
	// ROC             TierStage
	// PreferredReturn TierStage
	// Catchup         TierStage
	// FinalSplit      TierStage
}

type Contribution struct {
	Date              time.Time
	Amount            *money.Money
	ReturnCapitalLeft *money.Money
}

type TierStageType string

const (
	RocTierStage         TierStageType = "roc"
	PreferredReturnStage TierStageType = "preferred-return"
	CatchUpStage         TierStageType = "catchup-stage"
	FinalSplitStage      TierStageType = "final-split"
)

type TierStage struct {
	TierName          TierStageType
	StartingCapital   *money.Money
	LpAllocattion     *money.Money
	GpAllocattion     *money.Money
	TotalDistribution *money.Money
	RemainingCapital  *money.Money
}

func (t TierStage) Display() string {
	return fmt.Sprintf("{TierName:%s, StartingCapital:%s, LpAllocattion:%s, GpAllocattion:%s, TotalDistribution:%s, RemainingCapital:%s}",
		t.TierName,
		t.StartingCapital.Display(),
		t.LpAllocattion.Display(),
		t.GpAllocattion.Display(),
		t.TotalDistribution.Display(),
		t.RemainingCapital.Display(),
	)
}

func CalculateDistribution(cfg Config, transactions []*model.Transaction) error {
	state := State{}
	for _, t := range transactions {
		switch t.Operation {
		case model.ContributionOperation:
			contribution := Contribution{
				Date:              t.TransactionDate.Time,
				Amount:            money.NewFromFloat(t.Amount*-1, money.USD),
				ReturnCapitalLeft: money.NewFromFloat(t.Amount*-1, money.USD),
			}
			state.LpContributions = append(state.LpContributions, contribution)
		case model.DistributionOperation:
			fmt.Printf("\nDISTRIBUTION%+v\n", t)
			capital := money.NewFromFloat(t.Amount, money.USD)
			ts, contributions, err := calculateReturnOfCapital(capital, state.LpContributions)
			if err != nil {
				return fmt.Errorf("failed to calculate return of capital: %w", err)
			}
			state.LpContributions = contributions
			fmt.Printf("STAGE %+v\n", ts.Display())

			ts, err = calculatePreferredReturn(cfg.PreferedReturn, t.TransactionDate.Time, ts.RemainingCapital, contributions)
			if err != nil {
				return fmt.Errorf("failed to calculate preferred return: %w", err)
			}
			fmt.Printf("STAGE %+v\n", ts.Display())

			ts, err = calculateCatchUp(cfg.CatchUp, ts.RemainingCapital, ts.LpAllocattion)
			if err != nil {
				return fmt.Errorf("failed to calculate catchup: %w", err)
			}
			fmt.Printf("STAGE %+v\n", ts.Display())

		}
	}
	return nil
}

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

var ErrRocLeft = fmt.Errorf("contribution still has roc left")

func calculatePreferredReturn(cfg PreferedReturnConfig, date time.Time, startingCapital *money.Money, contributions []Contribution) (TierStage, error) {
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
		rawPreferedReturn := cashFlow*math.Pow((1+cfg.HurdlePercentage), exponent) - cashFlow
		preferredReturn := money.NewFromFloat(rawPreferedReturn, money.USD)

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

func calculateCatchUp(cfg CatchUpConfig, startingCapital *money.Money, preferredReturn *money.Money) (TierStage, error) {
	cu := cfg.CatchupPercentage
	ci := cfg.CariedInterestPercentage
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

	distribution, err := startingCapital.Subtract(capitalLeft)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not subtract startingCapital and capitalLeft: %w", err)
	}
	remainingCapital, err := startingCapital.Subtract(distribution)
	if err != nil {
		return TierStage{}, fmt.Errorf("could not subtract startingCapital and distribution: %w", err)
	}
	ret := TierStage{
		TierName:          CatchUpStage,
		StartingCapital:   startingCapital,
		LpAllocattion:     money.New(0, money.USD),
		GpAllocattion:     distribution,
		TotalDistribution: distribution,
		RemainingCapital:  remainingCapital,
	}

	return ret, nil
}

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
