package cmd

import (
	"github.com/go-playground/redis/samples"
	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use: "lock",
	Run: func(cmd *cobra.Command, args []string) {
		samples.LockMain()
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
}
