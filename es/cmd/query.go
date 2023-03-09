package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
	"log"
	"time"
)

const (
	indexName = "test-index"
)

// https://www.elastic.co/blog/the-go-client-for-elasticsearch-introduction

var queryCmd = &cobra.Command{
	Use: "query",
	Run: func(cmd *cobra.Command, args []string) {
		query()
	},
}

func query() {
	es, _ := elasticsearch.NewDefaultClient()
	res, _ := es.Indices.Create(indexName)
	log.Printf("creating index result %v", res)

	fin1 := make(chan int)
	fin2 := make(chan int)
	docs := []string{
		"Document 1",
		"Wakanda Forever",
		"Super man",
		"Iron man",
	}
	docLen := len(docs)
	go periodicIndexing(es, fin1, docs)
	go periodicQuery(es, fin2, docs, docLen)

	<-fin1
	log.Println("Indexing finished")
	<-fin2
	log.Println("Querying finished")
}

// periodicIndexing to index document to es
func periodicIndexing(es *elasticsearch.Client, fin chan<- int, docs []string) {
	log.Println("begin indexing....")
	for id, doc := range docs {
		log.Printf("Indexing doc=%s ...", doc)
		//document := struct {
		//	Id   int    `json:"id"`
		//	Name string `json:"name"`
		//}{
		//	Id:   id,
		//	Name: doc,
		//}
		//es.Index(indexName, document)
		time.Sleep(100 * time.Millisecond)
	}
	fin <- 1

}

func periodicQuery(es *elasticsearch.Client, fin chan<- int, docs []string, docLen int) {
	log.Printf("begin querying %v docs...", docLen)
	queriedDoc := 0

	for queriedDoc < docLen {
		log.Printf("queriedDoc: %v query=%v", queriedDoc, docs[queriedDoc])
		queriedDoc += 1
		time.Sleep(100 * time.Millisecond)
	}
	fin <- 1
}
