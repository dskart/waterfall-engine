package cmd

import (
	"os"
	"path/filepath"

	"github.com/dskart/netword/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func init() {
	rootLogger, _ = logger.NewLogger(false)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

var rootLogger *zap.Logger

var rootCmd = &cobra.Command{
	Use:           filepath.Base(os.Args[0]),
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rootLogger.Fatal(err.Error())
	}
}
