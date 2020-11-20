package main

import (
	json2 "encoding/json"
	"fmt"
	"log"
	"strings"
	"syscall"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/crypto/ssh/terminal"
)

// MeiliSearch will index the sentences
// on a given MeiliSearch instance.
type MeiliSearch struct {
	client                 meilisearch.ClientInterface
	host, APIKey, indexUID string
	APIKeyRequired         bool
}

// Init the MeiliSearch client and the index.
func (m *MeiliSearch) Init() {
	// Format the host.
	var host string

	if strings.HasPrefix(m.host, "http://") || strings.HasPrefix(m.host, "https://") {
		host = m.host
	} else {
		host = "http://" + m.host
	}

	// Ask the API key if needed.
	if m.APIKeyRequired {
		m.askAPIKey()
	}

	// Set the index name.
	m.indexUID = "tatoeba"

	// Create a MeiliSearch client.
	m.client = meilisearch.NewClient(meilisearch.Config{
		Host:   host,
		APIKey: m.APIKey,
	})

	// Check if the index exist.
	if _, err := m.client.Indexes().Get(m.indexUID); err != nil {
		// The index doesn't exist, create it.
		m.createIndex()
	}

	// Set the searchable attributes.
	m.setSearchableAttributes()
}

// createIndex will create the Tatoeba index for Meilisearch.
func (m MeiliSearch) createIndex() {
	// Create an index if the index does not exist.
	_, err := m.client.Indexes().Create(meilisearch.CreateIndexRequest{
		Name: strings.Title(m.indexUID),
		UID:  m.indexUID,
	})

	if err != nil {
		log.Fatal(err)
	}
}

// setSearchableAttributes will set the searchable attributes.
func (m MeiliSearch) setSearchableAttributes() {
	searchableAttributes := []string{"id", "language", "content", "username"}

	m.client.Settings(m.indexUID).UpdateSearchableAttributes(searchableAttributes)
}

// Index sentences to the MeiliSearch instance.
func (m MeiliSearch) Index(sentences map[string]Sentence) {
	// Store the total of sentences.
	totalSentences := len(sentences)

	// Get the index documents from the client.
	index := m.client.Documents(m.indexUID)

	// i represent the current index of the loop.
	i := 1

	// Create a map of interfaces to store the sentences to index.
	var documents []map[string]interface{}

	// Loop over all sentences and index them.
	for _, sentence := range sentences {
		// Create an empty interface to convert the struct in.
		var sentenceInterface map[string]interface{}

		// Create a JSON from the struct to be able to
		// convert it into interface.
		sentenceAsJSON, _ := json2.Marshal(sentence)

		// Convert the JSON to the interface.
		json2.Unmarshal(sentenceAsJSON, &sentenceInterface)

		// Add the sentence interface to documents to index.
		documents = append(documents, sentenceInterface)

		// Call the API to add the sentences every 10000 sentences.
		if len(documents) == 10000 || i == totalSentences {
			// Add documents.
			addResponse, err := index.AddOrReplace(documents)

			if err != nil {
				log.Fatal(err)
			}

			// Reset the sentences array of map.
			documents = make([]map[string]interface{}, 0)

			// Log to the terminal the advance.
			fmt.Printf("\rIndexing sentences %d of %d", i, totalSentences)

			// Wait until the documents has been added by calling
			// the update API.
			for {
				// Wait 2 secondes between every update call.
				time.Sleep(2 * time.Second)

				// Get the update reponse.
				response, _ := m.client.Updates(m.indexUID).Get(addResponse.UpdateID)

				// Continue the indexation when the last update has been processed.
				if response.Status == meilisearch.UpdateStatusProcessed {
					break
				}
			}
		}

		// Increment the counter.
		i++
	}

	fmt.Println()
}

// askAPIKey will prompt in terminal to enter the API key.
func (m *MeiliSearch) askAPIKey() {
	// Ask user to enter the api key.
	fmt.Print("Please enter the API key: ")

	// Ask the user to enter the API key from the terminal.
	APIKey, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		log.Fatal(err)
	}

	// Store the API key.
	m.APIKey = string(APIKey)

	fmt.Println()
}
