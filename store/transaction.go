package store

import (
	"strconv"
	"time"

	"github.com/dskart/waterfall-engine/model"
	"github.com/google/uuid"
)

const (
	transactionKey                     = "transaction"
	transactionsByCommitmentIdKey      = "transactions_by_commitment_id"
	transactionsByCommitmentIdAndTsKey = "transactions_by_commitment_id_and_ts"
)

func (s *Store) AddTransaction(transaction model.Transaction) error {
	serialized, err := NewGzipSerializer().Serialize(transaction)
	if err != nil {
		return err
	}

	tx := s.backend.AtomicWrite()
	tx.SAdd(transactionsByCommitmentIdKey+":"+strconv.Itoa(transaction.CommitmentId), transaction.Id.String())
	tx.ZAdd(transactionsByCommitmentIdAndTsKey+":"+strconv.Itoa(transaction.CommitmentId), transaction.Id.String(), timeMicrosecondScore(transaction.TransactionDate.Time))
	tx.SetNX(transactionKey+":"+transaction.Id.String(), serialized)
	if ok, err := execAtomicWrite(tx); err != nil {
		return err
	} else if !ok {
		return ErrContention
	}
	return err
}

func (s *Store) GetTransactionsByIds(ids ...uuid.UUID) ([]*model.Transaction, error) {
	idStrings := make([]string, len(ids))
	for i, id := range ids {
		idStrings[i] = id.String()
	}

	var ret []*model.Transaction
	return ret, s.getByIds(transactionKey, &ret, NewGzipSerializer(), idStrings...)
}

func (s *Store) GetTransactionsByCommitmentId(commitmentId int) ([]*model.Transaction, error) {
	ids, err := s.backend.SMembers(transactionsByCommitmentIdKey + ":" + strconv.Itoa(commitmentId))
	if err != nil {
		return nil, err
	}
	return s.GetTransactionsByIds(stringsToUUIDs(ids)...)
}

// Gets transactions within an inclusive time range. If limit is non-zero, the returned events will be
// limited to that number. If limit is negative, the returned events will be the last events in the
// range.
func (s *Store) GetTransactionsByCommitmentIdAndTimeRange(commitmentId int, minTime, maxTime time.Time, limit int) ([]*model.Transaction, error) {
	ids, err := s.getIdsByKeyAndTimeRange(transactionsByCommitmentIdAndTsKey+":"+strconv.Itoa(commitmentId), minTime, maxTime, limit)
	if err != nil {
		return nil, err
	}
	return s.GetTransactionsByIds(stringsToUUIDs(ids)...)
}
