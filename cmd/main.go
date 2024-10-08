package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tonghia/transaction-history/internal/args"
	"github.com/tonghia/transaction-history/internal/processor"
)

func main() {
	// Set up logging to include the timestamp.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Define and parse command-line flags.
	interactivePtr := flag.Bool("interactive", false, "Enable interactive mode to input period and file path")
	periodPtr := flag.String("period", "", "Year and Month in YYYYMM format (required if not in interactive mode)")
	filePathPtr := flag.String("file", "", "Path to the CSV file containing transactions (required if not in interactive mode)")
	workernumPtr := flag.Int("workernum", 0, "Enable split file into chunk and process")
	outPathPtr := flag.String("out", "", "Path to the output file containing summary result in JSON format (optional)")

	flag.Parse()

	if *interactivePtr {
		// Run interactive input
		interactiveInput(periodPtr, filePathPtr)
	}

	// Parse command-line arguments
	yearMonth, err := args.ParsePeriod(*periodPtr)
	if err != nil {
		log.Fatalf("Invalid period: %v", err)
	}

	filePath, err := args.ParseFilePath(*filePathPtr)
	if err != nil {
		log.Fatalf("Invalid file path: %v", err)
	}

	summaryJSON, err := processor.Process(filePath, yearMonth, *workernumPtr)
	if err != nil {
		log.Fatalf("Error processing CSV file: %v", err)
	}

	if *outPathPtr != "" {
		if err := generateOutputFile(summaryJSON, *outPathPtr); err != nil {
			log.Fatalf("Error generating output file: %v", err)
		}
	} else {
		fmt.Println(string(summaryJSON))
	}
}

func interactiveInput(periodPtr, filePathPtr *string) {
	// Interactive Mode: Prompt the user for period and file path.
	fmt.Println("Interactive Mode Enabled.")
	reader := bufio.NewReader(os.Stdin)

	// Prompt for Period
	for {
		fmt.Print("Enter the period (YYYYMM): ")
		inputPeriod, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		inputPeriod = strings.TrimSpace(inputPeriod)
		if inputPeriod == "" {
			fmt.Println("Period cannot be empty. Please try again.")
			continue
		}

		*periodPtr = inputPeriod
		break
	}

	// Prompt for File Path
	for {
		fmt.Print("Enter the path to the CSV file: ")
		inputFilePath, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error reading input: %v", err)
		}
		inputFilePath = strings.TrimSpace(inputFilePath)
		if inputFilePath == "" {
			fmt.Println("File path cannot be empty. Please try again.")
			continue
		}

		*filePathPtr = inputFilePath
		break
	}
}

func generateOutputFile(jsonData json.RawMessage, outputPath string) error {
	// Create or truncate the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Write JSON data to the file
	_, err = outFile.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}
