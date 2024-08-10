package engine

import (
	"fmt"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/dskart/waterfall-engine/model"
)

type Engine struct {
	cfg Config
}

func NewEngine(cfg Config) Engine {
	return Engine{
		cfg: cfg,
	}
}

type Contribution struct {
	Date              time.Time
	Amount            *money.Money
	ReturnCapitalLeft *money.Money
}

// ComputeDistributions computes all distributions for a specific LP using a waterfall method
// `transactions` should be a date ascending slice of transactions (contributions and distributions) of a specific LP
// ComputeDistributions will return the allocation of each waterfall tier for each distribution
// The waterfall parameters are controled by the Config.
func (e Engine) ComputeDistributions(transactions []*model.Transaction) ([]Distribution, error) {
	lpContributions := []Contribution{}
	distributions := []Distribution{}
	for _, t := range transactions {
		switch t.Operation {
		case model.ContributionOperation:
			contribution := Contribution{
				Date:              t.TransactionDate.Time,
				Amount:            money.NewFromFloat(t.Amount*-1, money.USD),
				ReturnCapitalLeft: money.NewFromFloat(t.Amount*-1, money.USD),
			}
			lpContributions = append(lpContributions, contribution)
		case model.DistributionOperation:
			capital := money.NewFromFloat(t.Amount, money.USD)
			newDistribution := NewDistribution(t.TransactionDate.Time, capital)

			// ----------------- ROC -----------------
			ts, contributions, err := calculateReturnOfCapital(capital, lpContributions)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate return of capital: %w", err)
			}
			// we need to update the lpContributions so we can keep track of contributions that have already been paid back
			lpContributions = contributions
			newDistribution.ROC = ts

			// ----------------- Preferred Return -----------------
			ts, err = calculatePreferredReturn(e.cfg.PreferredReturn, t.TransactionDate.Time, ts.RemainingCapital, contributions)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate preferred return: %w", err)
			}
			newDistribution.PreferredReturn = ts

			// ----------------- CatchUp -----------------
			ts, err = calculateCatchUp(e.cfg.CatchUp, ts.RemainingCapital, ts.LpAllocattion)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate catchup: %w", err)
			}
			newDistribution.Catchup = ts

			// ----------------- Final Split -----------------
			ts, err = calculateFinalSplit(e.cfg.FinalSplit, ts.RemainingCapital)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate catchup: %w", err)
			}
			newDistribution.FinalSplit = ts

			distributions = append(distributions, newDistribution)
		}
	}

	return distributions, nil

}
