package main

import (
	json2 "encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/term"
)

// MeiliSearch will index the sentences
// on a given MeiliSearch instance.
type MeiliSearch struct {
	client         meilisearch.ClientInterface
	host, APIKey   string
	APIKeyRequired bool
}

// totalSentencesToIndexByRow define the total of sentences
// to index each bulk.
const totalSentencesToIndexByRow = 10000

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

	// Create a MeiliSearch client.
	m.client = meilisearch.NewClient(meilisearch.Config{
		Host:   host,
		APIKey: m.APIKey,
	})

	// Check if the index exist.
	if _, err := m.client.Indexes().Get(IndexName); err != nil {
		// The index doesn't exist, create it.
		m.createIndex()
	}

	// Set the searchable attributes.
	m.setSearchableAttributes()

	// Print the current instance.
	fmt.Printf("Indexing on MeiliSearch on the host \"%s\".\n", host)
}

// createIndex will create the Tatoeba index for Meilisearch.
func (m MeiliSearch) createIndex() {
	// Create an index if the index does not exist.
	_, err := m.client.Indexes().Create(meilisearch.CreateIndexRequest{
		Name: strings.Title(IndexName),
		UID:  IndexName,
	})

	if err != nil {
		log.Fatal(err)
	}
}

// setSearchableAttributes will set the searchable attributes.
func (m MeiliSearch) setSearchableAttributes() {
	searchableAttributes := []string{"id", "language", "content", "username"}

	m.client.Settings(IndexName).UpdateSearchableAttributes(searchableAttributes)
}

// Index sentences to the MeiliSearch instance.
func (m MeiliSearch) Index(sentences map[string]Sentence) {
	// Store the total of sentences.
	totalSentences := len(sentences)

	// Get the index documents from the client.
	index := m.client.Documents(IndexName)

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

		// Call the API to add the sentences.
		if len(documents) == totalSentencesToIndexByRow || i == totalSentences {
			// Check if the client still working.
			if err := m.client.Health().Get(); err != nil {
				color.Red("\nThe server isn't responding anymore... Can't index sentences...")
				os.Exit(0)
			}

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
				response, _ := m.client.Updates(IndexName).Get(addResponse.UpdateID)

				// Continue the indexation when the last update has been processed.
				if response.Status == meilisearch.UpdateStatusProcessed {
					break
				}
			}
		}

		// Increment the counter.
		i++
	}
}

// askAPIKey will prompt in terminal to enter the API key.
func (m *MeiliSearch) askAPIKey() {
	// Ask user to enter the api key.
	fmt.Print("Please enter the API key: ")

	// Ask the user to enter the API key from the terminal.
	APIKey, err := term.ReadPassword(int(syscall.Stdin))

	if err != nil {
		log.Fatal(err)
	}

	// Store the API key.
	m.APIKey = string(APIKey)

	fmt.Println()
}
