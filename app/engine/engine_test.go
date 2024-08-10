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
			assert.Equal(t, tc.expectedTierStage.TierName, ts.TierName)
			assert.Equal(t, tc.expectedTierStage.StartingCapital, ts.StartingCapital)
			assert.Equal(t, tc.expectedTierStage.LpAllocattion, ts.LpAllocattion)
			assert.Equal(t, tc.expectedTierStage.TotalDistribution, ts.TotalDistribution)
			assert.Equal(t, tc.expectedTierStage.RemainingCapital, ts.RemainingCapital)

			for i, c := range contributions {
				expectedContribution := tc.expectedContributions[i]
				assert.Equal(t, expectedContribution.Amount, c.Amount)
				assert.Equal(t, expectedContribution.ReturnCapitalLeft, c.ReturnCapitalLeft)
			}
		})
	}
}

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
		{
			name:              "MoreThanAYer",
			startingCapital:   money.NewFromFloat(2000, money.USD),
			contributions:     []Contribution{{Date: now.AddDate(-1, 0, -200), Amount: money.NewFromFloat(1000, money.USD), ReturnCapitalLeft: money.New(0, money.USD)}},
			expectedTierStage: TierStage{TierName: PreferredReturnStage, StartingCapital: money.NewFromFloat(2000, money.USD), LpAllocattion: money.NewFromFloat(126.51, money.USD), GpAllocattion: money.New(0, money.USD), TotalDistribution: money.NewFromFloat(126.51, money.USD), RemainingCapital: money.NewFromFloat(1873.49, money.USD)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts, err := calculatePreferredReturn(cfg, now, tc.startingCapital, tc.contributions)
			assert.ErrorIs(t, err, tc.err)
			assert.Equal(t, tc.expectedTierStage.TierName, ts.TierName)
			assert.Equal(t, tc.expectedTierStage.StartingCapital, ts.StartingCapital)
			assert.Equal(t, tc.expectedTierStage.LpAllocattion, ts.LpAllocattion)
			assert.Equal(t, tc.expectedTierStage.TotalDistribution, ts.TotalDistribution)
			assert.Equal(t, tc.expectedTierStage.RemainingCapital, ts.RemainingCapital)
		})
	}
}

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
			assert.Equal(t, tc.expectedTierStage.TierName, ts.TierName)
			assert.Equal(t, tc.expectedTierStage.StartingCapital, ts.StartingCapital)
			assert.Equal(t, tc.expectedTierStage.LpAllocattion, ts.LpAllocattion)
			assert.Equal(t, tc.expectedTierStage.TotalDistribution, ts.TotalDistribution)
			assert.Equal(t, tc.expectedTierStage.RemainingCapital, ts.RemainingCapital)
		})
	}
}

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
			assert.Equal(t, tc.expectedTierStage.StartingCapital, ts.StartingCapital)
			assert.Equal(t, tc.expectedTierStage.LpAllocattion, ts.LpAllocattion)
			assert.Equal(t, tc.expectedTierStage.TotalDistribution, ts.TotalDistribution)
			assert.Equal(t, tc.expectedTierStage.RemainingCapital, ts.RemainingCapital)
		})
	}
}
