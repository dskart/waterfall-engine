package model

import "github.com/Rhymond/go-money"

type Stats struct {
	TotalProfit           *money.Money
	TotalDistribution     *money.Money
	TotalCommitment       *money.Money
	Contributed           int // rounded percentage
	ContributionRemaining *money.Money
}
