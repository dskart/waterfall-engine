package engine

import (
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

func TestEngine_CalculateReturnOfCapital(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name                  string
		startingCapital       *money.Money
		contributions         []Contribution
		expectedTierStage     TierStage
		expectedContributions []Contribution
		err                   error
	}{
		{
			name:                  "SingleContribution",
			startingCapital:       money.NewFromFloat(2000, money.USD),
			contributions:         []Contribution{{Date: now, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.NewFromFloat(1000, money.USD)}},
			expectedTierStage:     TierStage{TierName: RocTierStage, StartingCapital: money.NewFromFloat(2000, money.USD), LpAllocattion: money.NewFromFloat(1000, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(1000, money.USD), RemainingCapital: money.NewFromFloat(1000, money.USD)},
			expectedContributions: []Contribution{{Date: now, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}},
			err:                   nil,
		},
		{
			name:                  "MultipleContributions",
			startingCapital:       money.NewFromFloat(2000, money.USD),
			contributions:         []Contribution{{Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.NewFromFloat(500, money.USD)}, {Date: now, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.NewFromFloat(1000, money.USD)}},
			expectedTierStage:     TierStage{TierName: RocTierStage, StartingCapital: money.NewFromFloat(2000, money.USD), LpAllocattion: money.NewFromFloat(1500, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(1500, money.USD), RemainingCapital: money.NewFromFloat(500, money.USD)},
			expectedContributions: []Contribution{{Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}, {Date: now, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}},
			err:                   nil,
		},
		{
			name:                  "NoCapitalLeft",
			startingCapital:       money.NewFromFloat(700, money.USD),
			contributions:         []Contribution{{Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.NewFromFloat(500, money.USD)}, {Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.NewFromFloat(500, money.USD)}},
			expectedTierStage:     TierStage{TierName: RocTierStage, StartingCapital: money.NewFromFloat(700, money.USD), LpAllocattion: money.NewFromFloat(700, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(700, money.USD), RemainingCapital: money.New(0, money.USD)},
			expectedContributions: []Contribution{{Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}, {Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.NewFromFloat(300, money.USD)}},
			err:                   nil,
		},
		{
			name:                  "ContributionAlreadyReturned",
			startingCapital:       money.NewFromFloat(2000, money.USD),
			contributions:         []Contribution{{Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}, {Date: now, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.NewFromFloat(1000, money.USD)}},
			expectedTierStage:     TierStage{TierName: RocTierStage, StartingCapital: money.NewFromFloat(2000, money.USD), LpAllocattion: money.NewFromFloat(1000, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(1000, money.USD), RemainingCapital: money.NewFromFloat(1000, money.USD)},
			expectedContributions: []Contribution{{Date: now, Amount: money.NewFromFloat(500, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}, {Date: now, Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}},
			err:                   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, contributions, err := calculateReturnOfCapital(tc.startingCapital, tc.contributions)
			assert.ErrorIs(t, err, tc.err)
			assertTierStage(t, tc.expectedTierStage, ts)

			for i, c := range contributions {
				expectedContribution := tc.expectedContributions[i]
				ok, err := expectedContribution.ReturnCapitalLeft.Equals(c.ReturnCapitalLeft)
				assert.NoError(t, err)
				assert.True(t, ok)
				ok, err = expectedContribution.Amount.Equals(c.Amount)
				assert.NoError(t, err)
				assert.True(t, ok)
			}
		})
	}
}
