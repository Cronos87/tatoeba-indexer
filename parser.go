package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// readCSV read Tatoeba's CSV and returns a reader.
func readCSV(filename string) *bufio.Scanner {
	// Check if the file exists.
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		color.Red("\nThe file \"%s\" doesn't exist.\n", filename)
		os.Exit(0)
	}

	// Open the sentences file.
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	return bufio.NewScanner(file)
}

// ParseSentences will parse the file `sentences_detailed.csv`
// and returns a map of `Sentence`.
func ParseSentences() map[string]Sentence {
	// Open the sentences file.
	scanner := readCSV("sentences_detailed.csv")

	// Scan lines.
	scanner.Split(bufio.ScanLines)

	// Create an empty map of sentences.
	sentences := make(map[string]Sentence)

	// Loop over all lines and create a struct of Sentence.
	for scanner.Scan() {
		// Read and split the current line.
		line := strings.Split(scanner.Text(), "\t")

		// If the language code is not 3 characters,
		// ignore the line.
		if len(line[1]) < 3 {
			continue
		}

		// Convert the id from a string to an int.
		id, err := strconv.Atoi(line[0])

		if err != nil {
			log.Fatal(err)
		}

		// Convert the added date to nil if the csv value
		// is equal to \N.
		addedAt := line[4]

		if addedAt == "\\N" || addedAt == "0000-00-00 00:00:00" {
			addedAt = ""
		}

		// Do the same as before for the updated date.
		updatedAt := line[5]

		if updatedAt == "\\N" || updatedAt == "0000-00-00 00:00:00" {
			updatedAt = ""
		}

		// Create a sentence and append it to the
		// slice to returns.
		sentences[line[0]] = Sentence{
			ID:                  int32(id),
			Language:            line[1],
			Content:             line[2],
			Username:            line[3],
			AddedAt:             addedAt,
			UpdatedAt:           updatedAt,
			DirectRelations:     make([]int32, 0),
			IndirectRelations:   make([]int32, 0),
			TranslatedLanguages: make([]string, 0),
			HasAudio:            false,
		}
	}

	return sentences
}

// ParseSentencesLink will parse the file `links.csv`
// and add direct translations between sentences.
func ParseSentencesLink(sentences *map[string]Sentence) {
	// Open the links file.
	scanner := readCSV("links.csv")

	// Scan lines.
	scanner.Split(bufio.ScanLines)

	// Loop over all lines and create a struct of SentencesLink.
	for scanner.Scan() {
		// Read the current line.
		line := strings.Split(scanner.Text(), "\t")

		// Get the sentences from the map.
		fromSentence, fromIDExist := (*sentences)[line[0]]
		toSentence, toIDExist := (*sentences)[line[1]]

		// Insert the relation if both ids exists in the sentences map.
		if fromIDExist && toIDExist {
			// Add the relation.
			fromSentence.DirectRelations = append(fromSentence.DirectRelations, toSentence.ID)

			// Add the translated language if not present in the array.
			languageExist := false

			// Loop over the translated languages array to find if the
			// language already exist.
			for _, language := range fromSentence.TranslatedLanguages {
				if language == toSentence.Language {
					languageExist = true
					break
				}
			}

			// Add the language if not exist.
			if !languageExist {
				fromSentence.TranslatedLanguages = append(fromSentence.TranslatedLanguages, toSentence.Language)
			}

			// Update sentence.
			(*sentences)[line[0]] = fromSentence
		}
	}
}

// FindIndirectRelations add indirect translations between sentences.
func FindIndirectRelations(sentences *map[string]Sentence) {
	// Loop over all sentences.
	for ID, sentence := range *sentences {
		// If the sentence haven't direct relation, skip here.
		if len(sentence.DirectRelations) == 0 {
			continue
		}

		// Loop over all direct relations.
		for _, directSentenceID := range sentence.DirectRelations {
			// Get the sentence from the ID.
			directSentence := (*sentences)[strconv.Itoa(int(directSentenceID))]

			// Loop over all the direct relation sentence direct relations.
		DirectRelationLoop:
			for _, directDirectSentenceID := range directSentence.DirectRelations {
				// Continue to the next sentence if the direc direct sentence ID
				// is the same as the sentence ID.
				if directDirectSentenceID == sentence.ID {
					continue
				}

				// If the sentence haven't this sentence
				// in the indirect relations, add it.
				for _, indirectRelationSentenceID := range sentence.IndirectRelations {
					// If the direct direct sentence ID has been founded
					// in the indirect relations, stop here.
					if directDirectSentenceID == indirectRelationSentenceID {
						continue DirectRelationLoop
					}
				}

				// Add the direct direct relation to the indirect relation.
				sentence.IndirectRelations = append(sentence.IndirectRelations, directDirectSentenceID)

				// Add the language in translated languages if haven't.
				languageFound := false

				// Get the direct direct sentence.
				directDirectSentence := (*sentences)[strconv.Itoa(int(directDirectSentenceID))]

				for _, language := range sentence.TranslatedLanguages {
					if language == directDirectSentence.Language {
						languageFound = true
						break
					}
				}

				// Add the language if not founded.
				if !languageFound {
					sentence.TranslatedLanguages = append(sentence.TranslatedLanguages, directDirectSentence.Language)
				}
			}
		}

		// Update sentence.
		(*sentences)[ID] = sentence
	}
}

// ParseSentencesWithAudio will parse the file `sentences_with_audio.csv`
// and update the list of `Sentence` flagging the HasAudio property to true
// if the sentence id has been found in this file.
func ParseSentencesWithAudio(sentences *map[string]Sentence) {
	// Open the links file.
	scanner := readCSV("sentences_with_audio.csv")

	// Scan lines.
	scanner.Split(bufio.ScanLines)

	// Loop over all lines and create a struct of SentencesLink.
	for scanner.Scan() {
		// Read the current line.
		line := strings.Split(scanner.Text(), "\t")

		// Update the flag to true in the sentences map.
		var sentence = (*sentences)[line[0]]
		sentence.HasAudio = true
		(*sentences)[line[0]] = sentence
	}
}
