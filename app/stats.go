package app

import (
	"github.com/Rhymond/go-money"
	"github.com/dskart/waterfall-engine/model"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

func (s *Session) GetStatsByCommitmentId(commitmentId int) (model.Stats, appErrors.SanitizedError) {
	commitment, sanitizedErr := s.GetCommitmentById(commitmentId)
	if sanitizedErr != nil {
		return model.Stats{}, sanitizedErr
	}

	distributions, sanitizedErr := s.GetDistributionsByCommitmentId(commitmentId)
	if sanitizedErr != nil {
		return model.Stats{}, sanitizedErr
	}

	totalDistribution := money.New(0, money.USD)
	for _, d := range distributions {
		lpAllocation, err := d.LpTotalDistribution()
		if err != nil {
			return model.Stats{}, s.InternalError(err)
		}

		totalDistribution, err = totalDistribution.Add(lpAllocation)
		if err != nil {
			return model.Stats{}, s.InternalError(err)
		}
	}

	currentContribution, sanitizedErr := s.GetTotalContributionByCommitmentId(commitmentId)
	if sanitizedErr != nil {
		return model.Stats{}, sanitizedErr
	}

	totalCommitment := money.NewFromFloat(commitment.Amount, money.USD)
	totalProfit, err := totalDistribution.Subtract(currentContribution)
	if err != nil {
		return model.Stats{}, s.InternalError(err)
	}

	contributionRemaining, err := totalCommitment.Subtract(currentContribution)
	if err != nil {
		return model.Stats{}, s.InternalError(err)
	}
	contributed := currentContribution.AsMajorUnits() / totalCommitment.AsMajorUnits()

	return model.Stats{
		TotalProfit:           totalProfit,
		TotalDistribution:     totalDistribution,
		TotalCommitment:       totalCommitment,
		Contributed:           int(contributed * 100),
		ContributionRemaining: contributionRemaining,
	}, nil
}
