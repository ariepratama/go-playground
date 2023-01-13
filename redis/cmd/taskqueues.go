package cmd

import (
	"github.com/go-playground/redis/samples"
	"github.com/spf13/cobra"
)

var taskqueuesCmd = &cobra.Command{
	Use: "taskqueues",
	Run: func(cmd *cobra.Command, args []string) {
		samples.TaskQueueMain()
	},
}
