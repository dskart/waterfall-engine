package model

import (
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	ContributionOperation OperationType = "contribution"
	DistributionOperation OperationType = "distribution"
)

type Transaction struct {
	Id              uuid.UUID
	TransactionDate DateTime      `csv:"transaction_date"`
	Amount          float64       `csv:"transaction_amount"`
	Operation       OperationType `csv:"contribution_or_distribution"`
	CommitmentId    int           `csv:"commitment_id"`
}

type DateTime struct {
	time.Time
}

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format("02/01/2006"), nil
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("02/01/2006", csv)
	return err
}
