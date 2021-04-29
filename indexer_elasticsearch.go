package main

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

// Elasticsearch will index the sentences
// on a given Elasticsearch instance.
type Elasticsearch struct {
	client *elasticsearch.Client
	host   string
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

	// Declare the client init instance error.
	var err error

	// Create an Elasticsearch client.
	e.client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{host},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Print the current instance.
	fmt.Printf("Indexing on Elasticsearch on the host \"%s\".\n", host)
}

// Index sentences to the Elasticsearch instance.
func (e Elasticsearch) Index(sentences map[string]Sentence) {
	// Create a channel to receive the information a
	// sentence has been indexed.
	c := make(chan bool)

	// Store the total of sentences.
	totalSentences := len(sentences)

	// i represent the current index of the loop.
	i := 1

	// Loop over all sentences and index them.
	for index, sentence := range sentences {
		go func(sentence Sentence, index string, c chan bool) {
			// Create a JSON from the struct.
			sentenceAsJSON, _ := json2.Marshal(sentence)

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      "tatoeba",
				DocumentID: index,
				Body:       strings.NewReader(string(sentenceAsJSON)),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), e.client)

			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}

			defer res.Body.Close()

			c <- !res.IsError()
		}(sentence, index, c)

		// Get the answer from the channel if the sentence
		// has been indexed.
		isIndexed := <-c

		if !isIndexed {
			fmt.Println("non")
		}

		// Log to the terminal the advance.
		fmt.Printf("\rIndexing sentences %d of %d", i, totalSentences)

		// Increment the counter.
		i++
	}
}
