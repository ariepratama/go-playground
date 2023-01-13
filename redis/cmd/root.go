package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sample",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting redis sample")
	},
}

func init() {
	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(queuesCmd)
	rootCmd.AddCommand(taskqueuesCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
