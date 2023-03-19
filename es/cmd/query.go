package cmd

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/spf13/cobra"
	"log"
	"net/http"
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
	es, _ := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})

	cleanIndex(indexName)
	log.Println("Finished deleting index...")

	res, _ := es.Indices.Create(indexName).Do(context.Background())
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
func periodicIndexing(es *elasticsearch.TypedClient, fin chan<- int, docs []string) {
	log.Println("begin indexing....")
	for docId, doc := range docs {
		log.Printf("Indexing doc=%s ...", doc)
		document := struct {
			Id  int    `json:"id"`
			Doc string `json:"doc"`
		}{
			Id:  docId,
			Doc: doc,
		}
		resp, err := es.Index(indexName).
			Request(document).
			Do(context.Background())

		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Failed to create document to es: %v", err)
			return
		}

		log.Printf("successfully index document %v", document)

		//es.Index(indexName, http.Request(&es.Create.Request))
		time.Sleep(100 * time.Millisecond)
	}
	fin <- 1

}

func periodicQuery(es *elasticsearch.TypedClient, fin chan<- int, docs []string, docLen int) {
	log.Printf("begin querying %v docs...", docLen)
	queriedDoc := 0

	for queriedDoc < docLen {
		log.Printf("queriedDoc: %v query=%v", queriedDoc, docs[queriedDoc])

		searchRequest := &search.Request{
			Query: &types.Query{
				Match: map[string]types.MatchQuery{
					"doc": {Query: docs[queriedDoc]},
				},
			},
		}
		resp, err := es.Search().
			Index(indexName).
			Request(searchRequest).
			Do(context.Background())

		if err != nil || resp.StatusCode != http.StatusOK {
			log.Printf("Cannot seach document %v", err)

			time.Sleep(100 * time.Millisecond)
			continue
		}

		log.Printf("Successfully search a document %v", resp)
		queriedDoc += 1
		time.Sleep(100 * time.Millisecond)
	}
	fin <- 1
}
