package app

import (
	"time"

	"github.com/dskart/waterfall-engine/model"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

var (
	distantFuture time.Time = time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)
	distantPast   time.Time = time.Time{}
)

// GetSortedTransactionsByCommitmentId returns all transactions for a commitment id sorted by transaction date
func (s *Session) GetSortedTransactionsByCommitmentId(commitmentId int) ([]*model.Transaction, appErrors.SanitizedError) {
	minTime := distantPast
	transactions, err := s.store.GetTransactionsByCommitmentIdAndTimeRange(commitmentId, minTime, distantFuture, 0)
	if err != nil {
		return nil, s.InternalError(err)
	}
	return transactions, nil
}
