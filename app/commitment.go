package app

import (
	"fmt"
	"sort"

	"github.com/dskart/waterfall-engine/model"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

func (s *Session) GetCommitments() ([]*model.Commitment, appErrors.SanitizedError) {
	commitments, err := s.store.GetAllCommitments()
	if err != nil {
		return nil, s.InternalError(fmt.Errorf("could not get commitments: %w", err))
	}
	sort.Slice(commitments, func(i, j int) bool {
		return commitments[i].Id < commitments[j].Id
	})
	return commitments, nil
}

func (s *Session) GetCommitmentById(commitmentId int) (*model.Commitment, appErrors.SanitizedError) {
	commitments, err := s.store.GetCommitmentByIds(commitmentId)
	if err != nil {
		return nil, s.InternalError(fmt.Errorf("could not get commitment %d: %w", commitmentId, err))
	} else if len(commitments) != 1 {
		return nil, s.InternalError(fmt.Errorf("got more than 1 commitment for id %d", commitmentId))
	}
	return commitments[0], nil
}
