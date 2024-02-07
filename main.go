package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// latestGitTag will be set at build time using -ldflags
var latestGitTag string

// Function to retrieve the latest tag
func getLatestGitTag() string {
	return latestGitTag
}

func main() {
	showVersion := flag.Bool("v", false, "Prints the version")

	flag.Parse()

	if *showVersion {
		fmt.Println("Latest tag:", getLatestGitTag())
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Aborting. Needs two arguments: type of operation and input file.")
		os.Exit(1)
	}

	operationType := os.Args[1]
	src := os.Args[2]

	if _, err := os.Stat(src); os.IsNotExist(err) {
		fmt.Printf("Error: File '%s' not found.\n", src)
		os.Exit(1)
	}

	isFolder := false
	if info, err := os.Stat(src); err == nil && info.IsDir() {
		isFolder = true
	}

	if !isFolder {
		fmt.Printf("Processing file: %s\n", src)
		RunImageOperation(operationType, src)
	}

	if isFolder {
		fmt.Printf("Processing folder: %s\n", src)

		files, err := filepath.Glob(filepath.Join(src, "*"))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		}

		for _, file := range files {
			fmt.Printf("Found file: %s\n", file)
			if !IsFileImage(file) {
				fmt.Printf("Skipping non-image file: %s\n", file)
				continue
			}
			if IsFileAlreadyProcessed(file) {
				fmt.Printf("Skipping already processed file: %s\n", file)
				continue
			}
			fmt.Printf("Processing file: %s\n", file)
			RunImageOperation(operationType, file)
		}
	}
}
