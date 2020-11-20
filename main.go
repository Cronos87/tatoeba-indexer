package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
)

// Sentence describe the fields to index.
type Sentence struct {
	ID                  int32    `json:"id"`
	Language            string   `json:"language"`
	Content             string   `json:"content"`
	Username            string   `json:"username"`
	AddedAt             string   `json:"added_at"`
	UpdatedAt           string   `json:"updated_at"`
	DirectRelations     []int32  `json:"direct_translations"`
	IndirectRelations   []int32  `json:"indirect_translations"`
	TranslatedLanguages []string `json:"translated_languages"`
	HasAudio            bool     `json:"has_audio"`
}

func main() {
	// Declare arguments need to provide as CLI arguments.
	isAPIKeyRequired := flag.Bool("api-key", false, "will ask you to enter the API key")
	host := flag.String("host", "127.0.0.1:7700", "host url")

	// Parse CLI arguments.
	flag.Parse()

	// Create an instance of MeiliSearch.
	client := MeiliSearch{
		host:           *host,
		APIKeyRequired: *isAPIKeyRequired,
	}

	// Init the MeiliSearch client.
	client.Init()

	// Parse the sentences.
	fmt.Print("Parsing sentences...")
	sentences := ParseSentences()
	color.Green(fmt.Sprintf("%c[2K\rSentences has been parsed", 27))

	// Parse the audio file and update the sentences map.
	fmt.Print("Flag sentences with audio...")
	ParseSentencesWithAudio(&sentences)
	color.Green(fmt.Sprintf("%c[2K\rSentences with audio has been flagged", 27))

	// Parse the links between sentences and update the sentences map.
	fmt.Print("Add direct relations between sentences...")
	ParseSentencesLink(&sentences)
	color.Green(fmt.Sprintf("%c[2K\rDirect relations has been added", 27))

	// Add indirect relations between sentences.
	fmt.Print("Add indirect relations between sentences...")
	FindIndirectRelations(&sentences)
	color.Green(fmt.Sprintf("%c[2K\rIndirect relations has been added", 27))

	// Index sentences.
	client.Index(sentences)
}
