package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/dskart/waterfall-engine/app/engine"
	appErrors "github.com/dskart/waterfall-engine/pkg/errors"
)

func (s *Session) GetDistributionsByCommitmentId(commitmentId int) ([]engine.Distribution, appErrors.SanitizedError) {
	transactions, sanitizedErr := s.GetSortedTransactionsByCommitmentId(commitmentId)
	if sanitizedErr != nil {
		return nil, sanitizedErr
	}

	distributions, err := s.app.engine.ComputeDistributions(transactions)
	if err != nil {
		return nil, s.InternalError(fmt.Errorf("failed to calculate distribution: %w", err))
	}

	return distributions, nil
}

func (s *Session) GetWaterfallParameters() (engine.Config, appErrors.SanitizedError) {
	return s.app.engine.GetConfig(), nil
}

func (s *Session) GenerateAllDistributions() appErrors.SanitizedError {
	commitments, sanitizedErr := s.GetCommitments()
	if sanitizedErr != nil {
		return sanitizedErr
	}

	distributions := map[int][]engine.Distribution{}
	for _, c := range commitments {
		var sanitizedErr appErrors.SanitizedError
		distributions[c.Id], sanitizedErr = s.GetDistributionsByCommitmentId(c.Id)
		if sanitizedErr != nil {
			return sanitizedErr
		}
	}

	for id, d := range distributions {
		if err := saveDistributionsJsonFile("output", id, d); err != nil {
			return s.InternalError(err)
		}
	}

	return nil
}

func saveDistributionsJsonFile(dirPath string, commitmentId int, distributions []engine.Distribution) error {
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return err
	}

	location := path.Join(dirPath, fmt.Sprintf("%d.json", commitmentId))
	file, err := os.OpenFile(location, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	if err := enc.Encode(distributions); err != nil {
		return err
	}
	return nil
}
