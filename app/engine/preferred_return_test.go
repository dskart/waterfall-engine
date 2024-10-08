package engine

import (
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

func TestEngine_CalculatePreferredRetrun(t *testing.T) {
	cfg := PreferredReturnConfig{HurdlePercentage: 0.08}
	now := time.Now()
	yearAgo := now.AddDate(0, 0, -365)

	testCases := []struct {
		name              string
		startingCapital   *money.Money
		contributions     []Contribution
		expectedTierStage TierStage
		err               error
	}{
		{
			name:              "SingleContribution",
			startingCapital:   money.NewFromFloat(2000, money.USD),
			contributions:     []Contribution{{Date: yearAgo, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}},
			expectedTierStage: TierStage{TierName: PreferredReturnStage, StartingCapital: money.NewFromFloat(2000, money.USD), LpAllocattion: money.NewFromFloat(80, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(80, money.USD), RemainingCapital: money.NewFromFloat(1920, money.USD)},
			err:               nil,
		},
		{
			name:              "ErrRocLeft",
			startingCapital:   money.NewFromFloat(2000, money.USD),
			contributions:     []Contribution{{Date: yearAgo, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.NewFromFloat(1000, money.USD)}},
			expectedTierStage: TierStage{},
			err:               ErrRocLeft,
		},
		{
			name:              "MultipleContributions",
			startingCapital:   money.NewFromFloat(2000, money.USD),
			contributions:     []Contribution{{Date: yearAgo, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}, {Date: yearAgo, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}},
			expectedTierStage: TierStage{TierName: PreferredReturnStage, StartingCapital: money.NewFromFloat(2000, money.USD), LpAllocattion: money.NewFromFloat(120, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(120, money.USD), RemainingCapital: money.NewFromFloat(1880, money.USD)},
			err:               nil,
		},
		{
			name:              "NoCapitalLeft",
			startingCapital:   money.NewFromFloat(70, money.USD),
			contributions:     []Contribution{{Date: yearAgo, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}, {Date: yearAgo, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}},
			expectedTierStage: TierStage{TierName: PreferredReturnStage, StartingCapital: money.NewFromFloat(70, money.USD), LpAllocattion: money.NewFromFloat(70, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(70, money.USD), RemainingCapital: money.New(0, money.USD)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := calculatePreferredReturn(cfg, now, tc.startingCapital, tc.contributions)
			assert.ErrorIs(t, err, tc.err)
			if err == nil {
				assertTierStage(t, tc.expectedTierStage, ts)
			}
		})
	}
}
