package engine

import (
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/dskart/waterfall-engine/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func assertTierStage(t *testing.T, expected, actual TierStage) {
	assert.Equal(t, expected.TierName, actual.TierName)
	ok, err := expected.StartingCapital.Equals(actual.StartingCapital)
	assert.NoError(t, err)
	assert.True(t, ok)
	ok, err = expected.LpAllocattion.Equals(actual.LpAllocattion)
	assert.NoError(t, err)
	assert.True(t, ok)
	ok, err = expected.TotalDistribution.Equals(actual.TotalDistribution)
	assert.NoError(t, err)
	assert.True(t, ok)
	ok, err = expected.RemainingCapital.Equals(actual.RemainingCapital)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestEngine_ComputeDistributions(t *testing.T) {
	cfg := Config{
		PreferredReturn: PreferredReturnConfig{
			HurdlePercentage: 0.08,
		},
		CatchUp: CatchUpConfig{
			CatchupPercentage:         1.0,
			CarriedInterestPercentage: 0.2,
		},
		FinalSplit: FinalSplitConfig{
			LpPercentage: 0.8,
			GpPercentage: 0.2,
		},
	}

	engine := NewEngine(cfg)

	testCases := []struct {
		name                    string
		transactions            []*model.Transaction
		expectedROC             TierStage
		expectedPreferredReturn TierStage
		expectedCatchup         TierStage
		expectedFinalSplit      TierStage
	}{
		{
			name: "Example",
			transactions: []*model.Transaction{
				&model.Transaction{Id: uuid.New(), CommitmentId: 1, Operation: model.ContributionOperation, Amount: -1000.0, TransactionDate: model.DateTime{Time: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}},
				&model.Transaction{Id: uuid.New(), CommitmentId: 1, Operation: model.DistributionOperation, Amount: 2000.0, TransactionDate: model.DateTime{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}},
			},
			expectedROC: TierStage{
				TierName:          RocTierStage,
				StartingCapital:   money.NewFromFloat(2000, money.USD),
				LpAllocattion:     money.NewFromFloat(1000, money.USD),
				GpAllocattion:     money.NewFromFloat(0, money.USD),
				TotalDistribution: money.NewFromFloat(1000, money.USD),
				RemainingCapital:  money.NewFromFloat(1000, money.USD),
			},
			expectedPreferredReturn: TierStage{
				TierName:          PreferredReturnStage,
				StartingCapital:   money.NewFromFloat(1000, money.USD),
				LpAllocattion:     money.NewFromFloat(80, money.USD),
				GpAllocattion:     money.NewFromFloat(0, money.USD),
				TotalDistribution: money.NewFromFloat(80, money.USD),
				RemainingCapital:  money.NewFromFloat(920, money.USD),
			},
			expectedCatchup: TierStage{
				TierName:          CatchUpStage,
				StartingCapital:   money.NewFromFloat(920, money.USD),
				LpAllocattion:     money.NewFromFloat(0, money.USD),
				GpAllocattion:     money.NewFromFloat(20, money.USD),
				TotalDistribution: money.NewFromFloat(20, money.USD),
				RemainingCapital:  money.NewFromFloat(900, money.USD),
			},
			expectedFinalSplit: TierStage{
				TierName:          FinalSplitStage,
				StartingCapital:   money.NewFromFloat(900, money.USD),
				LpAllocattion:     money.NewFromFloat(720, money.USD),
				GpAllocattion:     money.NewFromFloat(180, money.USD),
				TotalDistribution: money.NewFromFloat(900, money.USD),
				RemainingCapital:  money.NewFromFloat(0, money.USD),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			distributions, err := engine.ComputeDistributions(tc.transactions)
			assert.NoError(t, err)
			assert.Len(t, distributions, 1)
			for _, d := range distributions {
				assertTierStage(t, tc.expectedROC, d.ROC)
				assertTierStage(t, tc.expectedPreferredReturn, d.PreferredReturn)
				assertTierStage(t, tc.expectedCatchup, d.Catchup)
				assertTierStage(t, tc.expectedFinalSplit, d.FinalSplit)
			}
		})
	}
}
