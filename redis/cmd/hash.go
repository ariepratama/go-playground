package cmd

import (
	"github.com/go-playground/redis/samples"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use: "hash",
	Run: func(cmd *cobra.Command, args []string) {
		samples.HashMain()
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
