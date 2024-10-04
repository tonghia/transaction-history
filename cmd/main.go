package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"transaction-history/pkg/transaction"
)

func main() {
	// Set up logging to include the timestamp.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse command-line arguments.
	yearMonth, filePath := parseArguments()

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	summary := transaction.ProcessData(file, yearMonth)

	// Generate JSON output.
	outputJSON(summary)
}

// parseArguments parses and validates command-line arguments.
func parseArguments() (string, string) {
	// Define command-line flags.
	yearMonthPtr := flag.String("period", "", "Year and Month in YYYYMM format (required)")
	filePathPtr := flag.String("file", "", "Path to the CSV file containing transactions (required)")

	// Parse the flags.
	flag.Parse()

	// Validate the inputs.
	if *yearMonthPtr == "" || *filePathPtr == "" {
		flag.Usage()
		log.Fatal("Both -period and -file arguments are required.")
	}

	// Verify that the file exists and is a regular file.
	absPath, err := filepath.Abs(*filePathPtr)
	if err != nil {
		log.Fatalf("Error resolving absolute path: %v", err)
	}

	fileInfo, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", absPath)
	}
	if err != nil {
		log.Fatalf("Error accessing file: %v", err)
	}
	if fileInfo.IsDir() {
		log.Fatalf("Provided file path is a directory, not a file: %s", absPath)
	}

	return *yearMonthPtr, absPath
}

// outputJSON marshals the summary into JSON and outputs it.
func outputJSON(summary transaction.Summary) {
	jsonData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
