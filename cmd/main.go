package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/tonghia/transaction-history/internal/args"
	"github.com/tonghia/transaction-history/internal/processor"
	"github.com/tonghia/transaction-history/internal/transaction"
)

type part struct {
	offset, size int64
}

func main() {
	// Set up logging to include the timestamp.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Define and parse command-line flags.
	interactivePtr := flag.Bool("interactive", false, "Enable interactive mode to input period and file path")
	periodPtr := flag.String("period", "", "Year and Month in YYYYMM format (required if not in interactive mode)")
	filePathPtr := flag.String("file", "", "Path to the CSV file containing transactions (required if not in interactive mode)")
	largeworkernumPtr := flag.Int("workernum", 0, "Enable split file into chunk and process")

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

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	var summary transaction.Summary
	chunkNum := *largeworkernumPtr
	if chunkNum <= 1 {
		summary, err = processor.ProcessData(file, yearMonth)
		if err != nil {
			log.Fatalf("Error processing CSV file: %v", err)
		}
	} else {
		// Determine non-overlapping parts for file split (each part has offset and size).
		parts, err := splitFile(filePath, chunkNum)
		if err != nil {
			log.Fatalf("Error spliting file")
		}
		// Start a goroutine to process each part, returning results on a channel.
		resultsCh := make(chan transaction.Summary)
		for _, part := range parts {
			go processPart(filePath, part.offset, part.size, yearMonth, resultsCh)
		}

		for i := 0; i < len(parts); i++ {
			result := <-resultsCh
			summary.TotalIncome = summary.TotalIncome + result.TotalIncome
			summary.TotalExpenditure = summary.TotalExpenditure + result.TotalExpenditure
			summary.Transactions = append(summary.Transactions, result.Transactions...)
		}

		summary.Period = yearMonth
		transaction.SortTransactions(summary.Transactions)
	}

	// Generate JSON output.
	outputJSON(summary)
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
}

// splitFile splits a file into multiple parts based on the specified number of parts
// refer to https://github.com/benhoyt/go-1brc/blob/fafba3256ea28631f6b3739f6d3b711a91199861/r8.go#L124
func splitFile(inputPath string, numParts int) ([]part, error) {
	const maxLineLength = 100 // Assumption
	const headerLength = 20

	f, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()
	splitSize := size / int64(numParts)

	buf := make([]byte, maxLineLength)

	parts := make([]part, 0, numParts)
	offset := int64(headerLength)
	for offset < size {
		seekOffset := max(offset+splitSize-maxLineLength, offset)
		if seekOffset > size {
			break
		}
		_, err := f.Seek(seekOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, _ := io.ReadFull(f, buf)
		chunk := buf[:n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			// Case 1: there is no `\n` at the end of the file, we stop
			// Case 2: maxLineLength is too small for the line, we accept there will be a huge chunk and improve it later
			break
		}
		remaining := len(chunk) - newline - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, part{offset, nextOffset - offset})
		offset = nextOffset
	}
	return parts, nil
}

func processPart(inputPath string, fileOffset int64, fileSize int64, yearMonth string, resultsCh chan<- transaction.Summary) {
	file, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Seek(fileOffset, io.SeekStart)
	if err != nil {
		panic(err)
	}
	// f := io.LimitedReader{R: file, N: fileSize}
	f := io.LimitReader(file, fileSize)

	summary, err := processor.ProcessLargeData(f, yearMonth)
	if err != nil {
		log.Fatalf("Error processing CSV file: %v", err)
	}

	resultsCh <- summary
}
