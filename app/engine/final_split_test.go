package engine

import (
	"testing"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

func TestEngine_CalculateFinalSplit(t *testing.T) {
	cfg := FinalSplitConfig{
		LpPercentage: 0.8,
		GpPercentage: 0.2,
	}

	testCases := []struct {
		name              string
		startingCapital   *money.Money
		expectedTierStage TierStage
		err               error
	}{
		{
			name:              "Split",
			startingCapital:   money.NewFromFloat(900, money.USD),
			expectedTierStage: TierStage{TierName: FinalSplitStage, StartingCapital: money.NewFromFloat(900.00, money.USD), LpAllocattion: money.NewFromFloat(720, money.USD), GpAllocattion: money.NewFromFloat(180, money.USD), TotalDistribution: money.NewFromFloat(900, money.USD), RemainingCapital: money.NewFromFloat(0, money.USD)},
			err:               nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := calculateFinalSplit(cfg, tc.startingCapital)
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.expectedTierStage.TierName, ts.TierName)
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
