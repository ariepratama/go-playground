package cmd

import (
	"github.com/go-playground/redis/samples"
	"github.com/spf13/cobra"
)

var aggresiveCacheCmd = &cobra.Command{
	Use: "agg",
	Run: func(cmd *cobra.Command, args []string) {
		samples.AggressiveCacheMain()
	},
}
