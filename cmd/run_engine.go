package cmd

import (
	"context"
	"fmt"

	"github.com/dskart/waterfall-engine/app"
	"github.com/dskart/waterfall-engine/pkg/shutdown"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runEngineCmd)
}

var runEngineCmd = &cobra.Command{
	Use:   "run-engine",
	Short: "load data and run waterfall engine",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		shutdown.OnShutdown(cancel)

		app, err := app.New(ctx, rootLogger, rootConfig.App)
		if err != nil {
			return fmt.Errorf("could not create app: %w", err)
		}

		session := app.NewSession(rootLogger).WithContext(ctx)
		session.Logger().Info("Running Engine...")
		if err := session.RunEngine(); err != nil {
			return err
		}
		session.Logger().Info("Done!")

		return nil
	},
}
