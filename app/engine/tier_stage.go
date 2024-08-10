package engine

import (
	"fmt"

	"github.com/Rhymond/go-money"
)

type TierStageType string

const (
	RocTierStage         TierStageType = "roc"
	PreferredReturnStage TierStageType = "preferred-return"
	CatchUpStage         TierStageType = "catchup-stage"
	FinalSplitStage      TierStageType = "final-split"
)

type TierStage struct {
	TierName          TierStageType
	StartingCapital   *money.Money
	LpAllocattion     *money.Money
	GpAllocattion     *money.Money
	TotalDistribution *money.Money
	RemainingCapital  *money.Money
}

func (t TierStage) Display() string {
	return fmt.Sprintf("{TierName:%s, StartingCapital:%s, LpAllocattion:%s, GpAllocattion:%s, TotalDistribution:%s, RemainingCapital:%s}",
		t.TierName,
		t.StartingCapital.Display(),
		t.LpAllocattion.Display(),
		t.GpAllocattion.Display(),
		t.TotalDistribution.Display(),
		t.RemainingCapital.Display(),
	)
}
