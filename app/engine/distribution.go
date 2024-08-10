package engine

import (
	"time"

	"github.com/Rhymond/go-money"
)

type Distribution struct {
	Date            time.Time
	Capital         *money.Money
	ROC             TierStage
	PreferredReturn TierStage
	Catchup         TierStage
	FinalSplit      TierStage
}

func NewDistribution(date time.Time, capital *money.Money) Distribution {
	return Distribution{
		Date:    date,
		Capital: capital,
	}
}

func (d Distribution) LpTotalDistribution() (*money.Money, error) {
	totalAllocation := money.New(0, money.USD)
	totalAllocation, err := totalAllocation.Add(d.ROC.LpAllocattion)
	if err != nil {
		return nil, err
	}

	totalAllocation, err = totalAllocation.Add(d.PreferredReturn.LpAllocattion)
	if err != nil {
		return nil, err
	}

	totalAllocation, err = totalAllocation.Add(d.Catchup.LpAllocattion)
	if err != nil {
		return nil, err
	}

	totalAllocation, err = totalAllocation.Add(d.FinalSplit.LpAllocattion)
	if err != nil {
		return nil, err
	}

	return totalAllocation, nil
}

func (d Distribution) DisplayLpTotalAllocation() (string, error) {
	lpAllocation, err := d.LpTotalDistribution()
	if err != nil {
		return "", nil
	}

	return lpAllocation.Display(), nil
}

func (d Distribution) DisplayGpTotalAllocation() (string, error) {
	totalAllocation := money.New(0, money.USD)
	totalAllocation, err := totalAllocation.Add(d.ROC.GpAllocattion)
	if err != nil {
		return "", err
	}

	totalAllocation, err = totalAllocation.Add(d.PreferredReturn.GpAllocattion)
	if err != nil {
		return "", err
	}

	totalAllocation, err = totalAllocation.Add(d.Catchup.GpAllocattion)
	if err != nil {
		return "", err
	}

	totalAllocation, err = totalAllocation.Add(d.FinalSplit.GpAllocattion)
	if err != nil {
		return "", err
	}

	return totalAllocation.Display(), nil
}

func (d Distribution) DisplayTotalDistribution() (string, error) {
	totalAllocation := money.New(0, money.USD)
	totalAllocation, err := totalAllocation.Add(d.ROC.TotalDistribution)
	if err != nil {
		return "", err
	}

	totalAllocation, err = totalAllocation.Add(d.PreferredReturn.TotalDistribution)
	if err != nil {
		return "", err
	}

	totalAllocation, err = totalAllocation.Add(d.Catchup.TotalDistribution)
	if err != nil {
		return "", err
	}

	totalAllocation, err = totalAllocation.Add(d.FinalSplit.TotalDistribution)
	if err != nil {
		return "", err
	}

	return totalAllocation.Display(), nil
}
