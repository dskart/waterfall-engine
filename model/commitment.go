package model

import "github.com/Rhymond/go-money"

type Commitment struct {
	Id         int     `csv:"id"`
	EntityName string  `csv:"entity_name"`
	FundId     int     `csv:"fund_id"`
	Amount     float64 `csv:"commitment_amount"`
}

func (c Commitment) DisplayAmount() string {
	ret := money.NewFromFloat(c.Amount, money.USD)
	return ret.Display()
}
