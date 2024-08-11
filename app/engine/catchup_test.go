package engine

import (
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

func TestEngine_CalculateCatchup(t *testing.T) {
	cfg := CatchUpConfig{
		CatchupPercentage:         1.0,
		CarriedInterestPercentage: 0.2,
	}

	testCases := []struct {
		name              string
		startingCapital   *money.Money
		preferedReturn    *money.Money
		expectedTierStage TierStage
		err               error
	}{
		{
			name:              "Catchup",
			startingCapital:   money.NewFromFloat(920.00, money.USD),
			preferedReturn:    money.NewFromFloat(80, money.USD),
			expectedTierStage: TierStage{TierName: CatchUpStage, StartingCapital: money.NewFromFloat(920.00, money.USD), LpAllocattion: money.New(0, money.USD), GpAllocattion: money.New(20, money.USD), TotalDistribution: money.NewFromFloat(20, money.USD), RemainingCapital: money.NewFromFloat(900, money.USD)},
			err:               nil,
		},
		{
			name:              "NoCapitalLeft",
			startingCapital:   money.NewFromFloat(10, money.USD),
			preferedReturn:    money.NewFromFloat(80, money.USD),
			expectedTierStage: TierStage{TierName: CatchUpStage, StartingCapital: money.NewFromFloat(10, money.USD), LpAllocattion: money.New(0, money.USD), GpAllocattion: money.New(10, money.USD), TotalDistribution: money.NewFromFloat(10, money.USD), RemainingCapital: money.NewFromFloat(0, money.USD)},
			err:               nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := calculateCatchUp(cfg, tc.startingCapital, tc.preferedReturn)
			assert.ErrorIs(t, err, tc.err)
			ok, err := tc.expectedTierStage.StartingCapital.Equals(ts.StartingCapital)
			assert.NoError(t, err)
			assert.True(t, ok)
			ok, err = tc.expectedTierStage.LpAllocattion.Equals(ts.LpAllocattion)
			assert.NoError(t, err)
			assert.True(t, ok)
			ok, err = tc.expectedTierStage.TotalDistribution.Equals(ts.TotalDistribution)
			assert.NoError(t, err)
			assert.True(t, ok)
			ok, err = tc.expectedTierStage.RemainingCapital.Equals(ts.RemainingCapital)
			assert.NoError(t, err)
			assert.True(t, ok)
		})
	}
}
