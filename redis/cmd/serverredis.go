package cmd

import (
	"github.com/go-playground/redis/servers"
	"github.com/spf13/cobra"
)

var serverredis1Cmd = &cobra.Command{
	Use: "serverredis1",
	Run: func(cmd *cobra.Command, args []string) {
		servers.InitServerRedis1()
	},
}

var serverredis2Cmd = &cobra.Command{
	Use: "serverredis2",
	Run: func(cmd *cobra.Command, args []string) {
		servers.InitServerRedis2()
	},
}
