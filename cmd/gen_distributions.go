package cmd

import (
	"context"
	"fmt"

	"github.com/dskart/waterfall-engine/app"
	"github.com/dskart/waterfall-engine/pkg/shutdown"
	"github.com/spf13/cobra"
)

func init() {
	genDistributionsCmd.Flags().StringP("data", "d", "./data", "the path to the data directory")
	rootCmd.AddCommand(genDistributionsCmd)
}

var genDistributionsCmd = &cobra.Command{
	Use:   "gen-distributions",
	Short: "generate distributions in json format",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		shutdown.OnShutdown(cancel)

		app, err := app.New(ctx, rootLogger, rootConfig.App)
		if err != nil {
			return fmt.Errorf("could not create app: %w", err)
		}

		dataPath, _ := cmd.Flags().GetString("data")
		session := app.NewSession(rootLogger).WithContext(ctx)
		session.Logger().Info("Loading data...")
		session.LoadData(dataPath)
		session.Logger().Info("Data loaded!")

		if err := session.GenerateAllDistributions(); err != nil {
			return err
		}

		return nil
	},
}
