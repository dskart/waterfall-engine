package app

import (
	"fmt"
	"os"

	"github.com/dskart/waterfall-engine/model"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
)

func (s *Session) LoadData() appErrors.SanitizedError {
	if err := s.loadCommitments(); err != nil {
		return s.InternalError(fmt.Errorf("could not load commiments: %w", err))
	}

	if err := s.loadTransactions(); err != nil {
		return s.InternalError(fmt.Errorf("could not load transactions: %w", err))
	}

	return nil
}

func (s *Session) loadCommitments() appErrors.SanitizedError {
	commitments := []*model.Commitment{}
	if err := loadFile("./data/commitments.csv", &commitments); err != nil {
		return s.InternalError(fmt.Errorf("could not load commitments: %w", err))
	}

	for _, c := range commitments {
		if err := s.store.AddCommitment(*c); err != nil {
			return s.InternalError(fmt.Errorf("could not add commitment %d: %w", c.Id, err))
		}
	}

	return nil
}

func (s *Session) loadTransactions() appErrors.SanitizedError {
	transactions := []*model.Transaction{}
	if err := loadFile("./data/transactions.csv", &transactions); err != nil {
		return s.InternalError(fmt.Errorf("could not load transactions: %w", err))
	}

	for _, t := range transactions {
		t.Id = uuid.New()
		if err := s.store.AddTransaction(*t); err != nil {
			return s.InternalError(fmt.Errorf("could not add transaction %d: %w", t.Id, err))
		}
	}

	return nil
}

func loadFile(location string, dest interface{}) error {
	file, err := os.Open(location)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", location, err)
	}
	defer file.Close()

	if err := gocsv.UnmarshalFile(file, dest); err != nil {
		return fmt.Errorf("could not unmarshal %s: %w", location, err)
	}
	return nil
}
