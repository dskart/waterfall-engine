package cmd

import (
	"context"

	"github.com/dskart/waterfall-engine/pkg/shutdown"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(noopCmd)
}

var noopCmd = &cobra.Command{
	Use:   "noop",
	Short: "does nothing",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, cancel := context.WithCancel(context.Background())
		shutdown.OnShutdown(cancel)
		rootLogger.Sugar().Infof("HELLO")
		return nil
	},
}
