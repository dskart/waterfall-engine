package model

type Commitment struct {
	Id         int     `csv:"id"`
	EntityName string  `csv:"entity_name"`
	FundId     int     `csv:"fund_id"`
	Amount     float64 `csv:"commitment_amount"`
}
