package main

import (
	"fmt"
	"log"
	"strings"

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
}

// Index sentences to the Elasticsearch instance.
func (m Elasticsearch) Index(sentences map[string]Sentence) {
	fmt.Println("Index Elasticsearch")
}
