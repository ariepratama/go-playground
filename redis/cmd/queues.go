package cmd

import (
	"github.com/go-playground/redis/samples"
	"github.com/spf13/cobra"
)

var queuesCmd = &cobra.Command{
	Use: "queues",
	Run: func(cmd *cobra.Command, args []string) {
		samples.QueueMain()
	},
}
