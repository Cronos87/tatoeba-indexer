package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/fatih/color"
	"github.com/mholt/archiver"
)

// DownloadFiles download the needed files to index Tatoeba's sentences.
func DownloadFiles(force bool) {
	// Create the array of files to download.
	files := []string{
		SentencesDetailed,
		Links,
		SentencesWithAudio,
	}

	// Log to the console that the files will be downloaded.
	if !force {
		fmt.Print(color.CyanString("Files doesn't exist."))
	} else {
		fmt.Print(color.CyanString("Force download asked."))
	}

	color.Cyan(" %d files need to be downloaded.", len(files))

	// Create client.
	client := grab.NewClient()

	// Download files and extract the content.
	for _, filename := range files {
		// Delete the CSV and tar.bz2 files.
		_ = os.Remove(fmt.Sprintf("%s%s.csv", os.TempDir(), filename))
		_ = os.Remove(fmt.Sprintf("%s%s.tar.bz2", os.TempDir(), filename))

		// Store the filename with the extension.
		filenameExt := filename + ".tar.bz2"

		// Format the filename to be easier to read.
		filenameFormatted := strings.Title(strings.ReplaceAll(filename, "_", " "))

		// Create the request.
		req, _ := grab.NewRequest(os.TempDir(), fmt.Sprintf("https://downloads.tatoeba.org/exports/%s", filenameExt))

		// Start the download.
		resp := client.Do(req)

		// Start UI loop.
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-t.C:
				// Log the progress.
				fmt.Printf("%s: downloaded %.2f%%\r", filenameFormatted, 100*resp.Progress())

			case <-resp.Done:
				// The download is complete, stop here.
				break Loop
			}
		}

		// Log the progress as downloaded.
		color.Green(fmt.Sprintf("%c[2K\r%s: downloaded", 27, filenameFormatted))

		// Check for errors.
		if err := resp.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			os.Exit(1)
		}

		// Log that the archive is beeing unarchiving.
		fmt.Printf("%s: unarchiving", filenameFormatted)

		// Extract the file.
		err := archiver.Unarchive(fmt.Sprintf("%s%s", os.TempDir(), filenameExt), os.TempDir())

		if err != nil {
			color.Red(fmt.Sprintf("Error while unarchiving the file %s.", filenameExt))
			os.Exit(1)
		}

		// Log that the archive has been unarchived.
		color.Green(fmt.Sprintf("%c[2K\r%s: unarchived", 27, filenameFormatted))
	}

	fmt.Println()
}
