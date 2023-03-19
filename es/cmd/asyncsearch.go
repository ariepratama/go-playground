package cmd

import (
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	asyncsearchSubmit "github.com/elastic/go-elasticsearch/v8/typedapi/asyncsearch/submit"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/spf13/cobra"
	"io"
	"log"
	"time"
)

const (
	asyncSearchIndexName = "async-search-index"
)

type (
	AsyncSearchResponse struct {
		Id        string                  `json:"id"`
		IsRunning bool                    `json:"is_running"`
		Response  AsyncSearchResponseResp `json:"response"`
	}

	AsyncSearchResponseResp struct {
		Hits AsyncSearchResponseHits `json:"hits"`
	}

	AsyncSearchResponseHits struct {
		Hits []AsyncSearchResponseHitsHits `json:"hits"`
	}

	AsyncSearchResponseHitsHits struct {
		Index  string                        `json:"_index"`
		Type   string                        `json:"_type"`
		Id     string                        `json:"_id"`
		Score  float32                       `json:"_score"`
		Source AsyncSearchResponseHitsSource `json:"_source"`
	}

	AsyncSearchResponseHitsSource struct {
		Id  int    `json:"id"`
		Doc string `json:"doc"`
	}

	AsyncSearchGetResponse struct {
		Id        string `json:"id"`
		IsRunning bool   `json:"is_running"`
	}
)

// https://www.elastic.co/guide/en/elasticsearch/reference/current/async-search.html
var asyncSearchCmd = &cobra.Command{
	Use: "asyncsearch",
	Run: func(cmd *cobra.Command, args []string) {
		asyncSearch()
	},
}

// asyncSearch simulate parallel indexing and searching asynchronously
func asyncSearch() {

	sampleDocs := []string{
		"document 1: Honey, This Mirror Isn't Big Enough For The Two Of Us",
		"document 2: I Never Told You What I Do For A Living",
		"document 3: I'm Not Okay (I Promise)",
		"document 3: It's Not A Fashion Statement, It's A Deathwish",
	}
	es, _ := elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})

	cleanIndex(asyncSearchIndexName)
	log.Println("finished removing index...")

	res, e := es.Indices.Create(asyncSearchIndexName).Do(context.Background())
	log.Printf("creating index result %v\n", res)
	log.Printf("creating index err %v\n", e)
	createIndexChan := asyncCreateIndex(es, sampleDocs)
	time.Sleep(100 * time.Millisecond)
	asyncSearchChan := asyncSearchIndex(es, sampleDocs)

	<-createIndexChan
	log.Printf("finished indexing all documents...")
	<-asyncSearchChan
	log.Printf("finished searching all documents...")
}

func asyncCreateIndex(es *elasticsearch.TypedClient, docs []string) <-chan int {
	createIndexChannel := make(chan int)
	go func() {
		for docId, doc := range docs {
			log.Printf("indexing document %v\n", doc)
			document := struct {
				Id  int    `json:"id"`
				Doc string `json:"doc"`
			}{
				Id:  docId,
				Doc: doc,
			}
			indexResp, indexErr := es.Index(asyncSearchIndexName).
				Request(document).
				Do(context.Background())
			log.Printf("response: %v", indexResp)

			if indexErr != nil {
				log.Printf("indexing error: %v", indexErr)
			}

			time.Sleep(100 * time.Millisecond)
		}

		createIndexChannel <- 1
	}()
	return createIndexChannel
}

func asyncSearchIndex(es *elasticsearch.TypedClient, docs []string) <-chan int {
	searchIndexChannel := make(chan int)

	go func() {
		for _, doc := range docs {
			log.Printf("searching document %v", doc)
			searchRequest := &asyncsearchSubmit.Request{
				Query: &types.Query{
					Match: map[string]types.MatchQuery{
						"doc": {Query: doc},
					},
				},
			}
			const maxRetry = 3
			var retryCount = 0

			// send async search request
			resp, err := es.Async.Submit().
				Request(searchRequest).
				Do(context.Background())

			// retry
			for err != nil && retryCount < maxRetry {
				log.Printf("error when asyncSearch: %v, retrying...", err)

				resp, err = es.Async.Submit().
					Request(searchRequest).
					Do(context.Background())
				retryCount++
				time.Sleep(200 * time.Millisecond)
			}

			bodyBytes, _ := io.ReadAll(resp.Body)
			log.Printf("response from asyncSearch: %v", resp)
			log.Printf("response from asyncSearch: %v", string(bodyBytes))

			var asyncSearchResp AsyncSearchResponse
			e := json.Unmarshal(bodyBytes, &asyncSearchResp)

			if e != nil {
				panic("cannot unmarshall response" + e.Error())
			}

			log.Printf("asyncSearch response body: %v\n", asyncSearchResp)
			log.Printf("asyncSearch response body: %v\n", string(bodyBytes))

			const maxPollingAttempt = 5
			var pollingCount = 0
			var isSearchFinished = asyncSearchResp.IsRunning

			for !isSearchFinished && !(pollingCount >= maxPollingAttempt) {
				var searchId = asyncSearchResp.Id

				getResp, _ := es.Async.Get(searchId).Do(context.Background())
				getRespBytes, _ := io.ReadAll(getResp.Body)
				json.Unmarshal(getRespBytes, &asyncSearchResp)
				time.Sleep(200 * time.Millisecond)
				pollingCount++
				isSearchFinished = asyncSearchResp.IsRunning
			}

			log.Printf("asyncSearch hits: %v", asyncSearchResp.Response.Hits)

		}
		searchIndexChannel <- 1
	}()
	return searchIndexChannel
}
