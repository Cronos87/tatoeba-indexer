package main

import (
	"bytes"
	"context"
	json2 "encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

// Elasticsearch will index the sentences
// on a given Elasticsearch instance.
type Elasticsearch struct {
	client      *elasticsearch.Client
	bulkIndexer esutil.BulkIndexer
	host        string
}

// Init the MeiliSearch client.
func (e *Elasticsearch) Init() {
	// Format the host.
	var host string

	if strings.HasPrefix(e.host, "http://") || strings.HasPrefix(e.host, "https://") {
		host = e.host
	} else {
		host = "http://" + e.host
	}

	// Declare the backoff function.
	retryBackoff := backoff.NewExponentialBackOff()

	// Declare the client init instance error.
	var err error

	// Create an Elasticsearch client.
	e.client, err = elasticsearch.NewClient(elasticsearch.Config{
		RetryOnStatus: []int{502, 503, 504, 429},
		Addresses:     []string{host},
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		MaxRetries: 5,
	})

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Delete the index
	res, err := e.client.Indices.Delete([]string{IndexName}, e.client.Indices.Delete.WithIgnoreUnavailable(true))

	if err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}

	res.Body.Close()

	// Re-create the index
	res, err = e.client.Indices.Create(IndexName)

	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}

	res.Body.Close()

	e.bulkIndexer, err = esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         IndexName,
		Client:        e.client,
		NumWorkers:    8,            // @TODO: Set this with a CLI variable.
		FlushBytes:    int(1000000), // @TODO: Set this with a CLI variable.
		FlushInterval: 30 * time.Second,
	})

	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	// Print the current instance.
	fmt.Printf("Indexing on Elasticsearch on the host \"%s\".\n", host)
}

// Index sentences to the Elasticsearch instance.
func (e Elasticsearch) Index(sentences map[string]Sentence) {
	// Store the total of sentences.
	totalSentences := len(sentences)

	// i represent the current index of the loop.
	i := 1

	// Loop over all sentences and index them.
	for index, sentence := range sentences {
		// Create a JSON from the struct.
		sentenceAsJSON, err := json2.Marshal(sentence)

		if err != nil {
			log.Fatalf("Cannot encode sentence %d: %s", sentence.ID, err)
		}

		// Add an item to the BulkIndexer
		err = e.bulkIndexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: index,
				Body:       bytes.NewReader(sentenceAsJSON),
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					// Log to the terminal the advance.
					fmt.Printf("\rIndexing sentences %d of %d", i, totalSentences)

					// Incremente the counter.
					i++
				},
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	// Close the indexer
	if err := e.bulkIndexer.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}
