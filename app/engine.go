package app

import (
	"fmt"

	"github.com/dskart/waterfall-engine/app/engine"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

func (s *Session) RunEngine() appErrors.SanitizedError {
	if sanitizedErr := s.LoadData(); sanitizedErr != nil {
		return sanitizedErr
	}

	transactions, sanitizedErr := s.GetSortedTransactionsByCommitmentId(1)
	if sanitizedErr != nil {
		return sanitizedErr
	}

	engineCfg := engine.Config{
		PreferedReturn: engine.PreferedReturnConfig{
			HurdlePercentage: 0.08,
		},
		CatchUp: engine.CatchUpConfig{
			CatchupPercentage:        1.0,
			CariedInterestPercentage: 0.2,
		},
		FinalSplit: engine.FinalSplitConfig{
			LpPercentage: 0.8,
			GpPercentage: 0.2,
		},
	}
	if err := engine.CalculateDistribution(engineCfg, transactions); err != nil {
		return s.InternalError(fmt.Errorf("failed to calculate distribution: %w", err))
	}

	return nil
}
