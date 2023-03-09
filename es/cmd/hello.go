package cmd

import "github.com/spf13/cobra"

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

var helloCmd = &cobra.Command{
	Use: "hello",
	Run: func(cmd *cobra.Command, args []string) {
		hello()
	},
}

// hello world to elastic search
// will fetch node information and display on the log
func hello() {
	es, _ := elasticsearch.NewDefaultClient()

	res, _ := es.Info()
	defer res.Body.Close()

	log.Println(res)
}
