package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/tonghia/transaction-history/pkg/transaction"
)

func main() {
	// Set up logging to include the timestamp.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Define and parse command-line flags.
	interactivePtr := flag.Bool("interactive", false, "Enable interactive mode to input period and file path")
	periodPtr := flag.String("period", "", "Year and Month in YYYYMM format (required if not in interactive mode)")
	filePathPtr := flag.String("file", "", "Path to the CSV file containing transactions (required if not in interactive mode)")

	flag.Parse()

	if *interactivePtr {
		// Run interactive input
		interactiveInput(periodPtr, filePathPtr)
	}

	// Parse command-line arguments
	yearMonth, filePath := parseArguments(periodPtr, filePathPtr)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	summary, err := transaction.ProcessData(file, yearMonth)
	if err != nil {
		log.Fatalf("Error processing CSV file: %v", err)
	}

	// Generate JSON output.
	outputJSON(summary)
}

// parseArguments parses and validates command-line arguments.
func parseArguments(yearMonthPtr, filePathPtr *string) (string, string) {

	// Validate the inputs.
	if *yearMonthPtr == "" || *filePathPtr == "" {
		flag.Usage()
		log.Fatal("Both -period and -file arguments are required.")
	}

	// Validate the period format.
	// _, _, err = transaction.ParseYearMonth(inputPeriod)
	// if err != nil {
	// 	fmt.Printf("Invalid period: %v. Please try again.\n", err)
	// 	continue
	// }

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

	return
}
