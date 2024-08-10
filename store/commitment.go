package store

import (
	"strconv"

	"github.com/dskart/waterfall-engine/model"
)

const (
	commitmentKey              = "commitment"
	commitmentsByEntityNameKey = "commitments_by_entity_name"
	allCommitmentsKey          = "all_commitments"
)

func (s *Store) AddCommitment(commitment model.Commitment) error {
	serialized, err := NewGzipSerializer().Serialize(commitment)
	if err != nil {
		return err
	}

	commitmentId := strconv.Itoa(commitment.Id)

	tx := s.backend.AtomicWrite()
	tx.SAdd(commitmentsByEntityNameKey+":"+commitment.EntityName, commitmentId)
	tx.SetNX(commitmentKey+":"+commitmentId, serialized)
	tx.SAdd(allCommitmentsKey, commitmentId)
	if ok, err := execAtomicWrite(tx); err != nil {
		return err
	} else if !ok {
		return ErrContention
	}
	return err
}

func (s *Store) GetCommitmentByIds(ids ...int) ([]*model.Commitment, error) {
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = strconv.Itoa(id)
	}
	var ret []*model.Commitment
	return ret, s.getByIds(commitmentKey, &ret, NewGzipSerializer(), idStrings...)
}

func (s *Store) GetAllCommitments() ([]*model.Commitment, error) {
	ids, err := s.backend.SMembers(allCommitmentsKey)
	if err != nil {
		return nil, err
	}
	var ret []*model.Commitment
	return ret, s.getByIds(commitmentKey, &ret, NewGzipSerializer(), ids...)
}
