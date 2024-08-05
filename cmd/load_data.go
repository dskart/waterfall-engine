package cmd

import (
	"context"
	"fmt"

	"github.com/dskart/waterfall-engine/app"
	"github.com/dskart/waterfall-engine/pkg/shutdown"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loadDataCmd)
}

var loadDataCmd = &cobra.Command{
	Use:   "load-data",
	Short: "load data from the ./data folder",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		shutdown.OnShutdown(cancel)

		app, err := app.New(ctx, rootLogger, rootConfig.App)
		if err != nil {
			return fmt.Errorf("could not create app: %w", err)
		}

		session := app.NewSession(rootLogger).WithContext(ctx)
		session.Logger().Info("Loading data...")
		session.LoadData()
		session.Logger().Info("Data loaded!")

		return nil
	},
}
