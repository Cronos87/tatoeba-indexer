package main

import (
	"fmt"
	"math"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/integrii/flaggy"
)

// Declare the files names.
const (
	SentencesDetailed  = "sentences_detailed"
	SentencesWithAudio = "sentences_with_audio"
	Links              = "links"
	Transcriptions     = "transcriptions"
)

// Declare the engines name.
const (
	meilisearchName   = "meilisearch"
	elasticsearchName = "elasticsearch"
)

// Declare CLI arguments variables and their defaults.
var IndexName = "tatoeba"
var needDownloadFiles = false

// MeiliSearch variables.
var isAPIKeyRequired = false
var hostMeiliSearch = "127.0.0.1:7700"

// Elasticsearch variables.
var hostElasticsearch = "127.0.0.1:9200"
var numWorkers = int(math.Min(2, float64(runtime.NumCPU())))
var flushBytes = 1000000

// parseCLIArguments will parse CLI arguments and populate
// the variables.
func parseCLIArguments() {
	// Create the global command.
	flaggy.String(&IndexName, "i", "index", "index name")
	flaggy.Bool(&needDownloadFiles, "d", "download-files", "download files needed to index Tatoeba's sentences")

	// Create the subcommand for MeiliSearch.
	meiliSearchSubcommand := flaggy.NewSubcommand(meilisearchName)
	meiliSearchSubcommand.Description = "Index sentences in MeiliSearch.\n\nhttps://www.meilisearch.com"

	// Declare arguments need to provide as CLI arguments.
	meiliSearchSubcommand.Bool(&isAPIKeyRequired, "", "api-key", "will ask you to enter the API key")
	meiliSearchSubcommand.String(&hostMeiliSearch, "", "host", "host url")

	// Create the subcommand for Elasticsearch.
	elasticsearchSubcommand := flaggy.NewSubcommand(elasticsearchName)
	elasticsearchSubcommand.Description = "Index sentences in Elasticsearch.\n\nhttps://www.elastic.co/elasticsearch/"

	// Declare arguments need to provide as CLI arguments.
	elasticsearchSubcommand.String(&hostElasticsearch, "", "host", "host url")
	elasticsearchSubcommand.Int(&numWorkers, "w", "workers", fmt.Sprintf("the number of workers. Maximum %d", runtime.NumCPU()))
	elasticsearchSubcommand.Int(&flushBytes, "b", "flush-bytes", "the flush threshold in bytes")

	// Add the subcommands to the parser.
	flaggy.AttachSubcommand(meiliSearchSubcommand, 1)
	flaggy.AttachSubcommand(elasticsearchSubcommand, 1)

	// Parse CLI arguments.
	flaggy.Parse()

	// Check if the number of workers is not exceeded.
	if numWorkers > runtime.NumCPU() {
		color.Cyan(fmt.Sprintf("You can't define more than %d workers. The value has been changed with the maximum one.", runtime.NumCPU()))
		numWorkers = runtime.NumCPU()
	}
}

// FileExists check if a file exists and returns true if exists, false otherwise.
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func main() {
	// Parse CLI arguments.
	parseCLIArguments()

	// If no subcommand was specified, show help and exit.
	if len(os.Args) == 1 {
		flaggy.ShowHelpAndExit("")
	}

	// Download files if needed.
	if needDownloadFiles ||
		!FileExists(os.TempDir()+SentencesDetailed+".csv") ||
		!FileExists(os.TempDir()+Links+".csv") ||
		!FileExists(os.TempDir()+SentencesWithAudio+".csv") {
		DownloadFiles(needDownloadFiles)
	}

	// Declare the client.
	var client Indexer

	// Create the client depending of the user choice.
	switch os.Args[1] {
	case meilisearchName:
		// Create an instance of MeiliSearch.
		client = &MeiliSearch{
			host:           hostMeiliSearch,
			APIKeyRequired: isAPIKeyRequired,
		}
	case elasticsearchName:
		// Create an instance of Elasticsearch.
		client = &Elasticsearch{
			host:       hostElasticsearch,
			numWorkers: numWorkers,
			flushBytes: flushBytes,
		}
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

	// Add some languages transcriptions.
	fmt.Print("Add transcriptions...")
	ParseTranscriptions(&sentences)
	color.Green(fmt.Sprintf("%c[2K\rTranscriptions has been added", 27))

	// Index sentences.
	client.Index(sentences)

	fmt.Println()
}
