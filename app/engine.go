package app

import (
	"fmt"

	"github.com/dskart/waterfall-engine/app/engine"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

func (s *Session) GetDistributionsByCommitmentId(commitmentId int) ([]engine.Distribution, appErrors.SanitizedError) {
	transactions, sanitizedErr := s.GetSortedTransactionsByCommitmentId(commitmentId)
	if sanitizedErr != nil {
		return nil, sanitizedErr
	}

	distributions, err := s.app.engine.ComputeDistributions(transactions)
	if err != nil {
		return nil, s.InternalError(fmt.Errorf("failed to calculate distribution: %w", err))
	}

	return distributions, nil
}

func (s *Session) GetWaterfallParameters() (engine.Config, appErrors.SanitizedError) {
	return s.app.engine.GetConfig(), nil
}
