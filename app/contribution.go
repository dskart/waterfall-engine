package app

import (
	"github.com/Rhymond/go-money"
	"github.com/dskart/waterfall-engine/model"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

func (s *Session) GetTotalContributionByCommitmentId(commitmentId int) (*money.Money, appErrors.SanitizedError) {
	transactions, sanitizedErr := s.GetSortedTransactionsByCommitmentId(commitmentId)
	if sanitizedErr != nil {
		return nil, sanitizedErr
	}

	totalAmount := 0.0
	for _, t := range transactions {
		if t.Operation == model.ContributionOperation {
			totalAmount += (-1 * t.Amount)
		}
	}

	ret := money.NewFromFloat(totalAmount, money.USD)

	return ret, nil
}
