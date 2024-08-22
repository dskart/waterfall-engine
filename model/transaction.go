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
	return date.Time.Format("01/02/2006"), nil
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("01/02/2006", csv)
	return err
}

func (date *DateTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len("01/02/2006")+2)
	b = append(b, '"')
	b = date.Time.UTC().AppendFormat(b, "01/02/2006")
	b = append(b, '"')
	return b, nil
}
