// csvgen.go
package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Record represents a single transaction.
type Record struct {
	Date    string
	Amount  int
	Content string
}

// Predefined lists for random data generation.
var (
	contents = []string{
		"Rent",
		"Salary",
		"Groceries",
		"Utilities",
		"Dining Out",
		"Entertainment",
		"Healthcare",
		"Transportation",
		"Investment",
		"Freelance Work",
		"Bonus",
		"Car Maintenance",
		"Subscriptions",
		"Shopping",
		"Books",
		"Gas",
		"Internet Bill",
		"Phone Bill",
		"Project Payment",
		"Debit",
		"Credit Card Payment",
		"Insurance",
		"Travel",
		"Education",
		"Miscellaneous",
	}
)

// Initialize the random seed.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// generateRandomDate generates a random date within the last 5 years.
func generateRandomDate() string {
	start := time.Now().AddDate(-5, 0, 0).Unix()
	end := time.Now().Unix()
	randomTime := time.Unix(rand.Int63n(end-start)+start, 0)
	return randomTime.Format("2006/01/02")
}

// generateRandomAmount generates a random amount.
// Positive values represent deposits, negative values represent withdrawals.
func generateRandomAmount() int {
	// Define the range for amounts.
	// Withdrawals: -1000 to -10
	// Deposits: 10 to 100000
	if rand.Intn(2) == 0 {
		return -rand.Intn(991) - 10 // -10 to -1000
	}
	return rand.Intn(99991) + 10 // 10 to 100000
}

// generateRandomContent selects a random content description.
func generateRandomContent() string {
	return contents[rand.Intn(len(contents))]
}

// validateInput checks if the input string is a valid uint64 number.
func validateInput(input string) (uint64, error) {
	num, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return 0, errors.New("invalid input: please provide a valid positive integer")
	}
	return num, nil
}

// generateRecords creates a slice of Records with random data.
func generateRecords(numRows uint64) []Record {
	records := make([]Record, numRows)
	for i := uint64(0); i < numRows; i++ {
		records[i] = Record{
			Date:    generateRandomDate(),
			Amount:  generateRandomAmount(),
			Content: generateRandomContent(),
		}
	}
	return records
}

// writeCSV writes the records to a CSV file with a header.
func writeCSV(records []Record, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"date", "amount", "content"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("error writing header to CSV: %v", err)
	}

	// Write records
	for _, record := range records {
		row := []string{
			record.Date,
			strconv.Itoa(record.Amount),
			record.Content,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing record to CSV: %v", err)
		}
	}

	return nil
}

func main() {
	// Define and parse command-line flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <number_of_rows>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Println("Error: Exactly one argument is required.")
		flag.Usage()
		os.Exit(1)
	}

	input := flag.Arg(0)

	// Validate input
	numRows, err := validateInput(input)
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	// Generate records
	fmt.Printf("Generating %d records...\n", numRows)
	records := generateRecords(numRows)

	// Define output CSV file name
	outputFile := "generated_transactions.csv"

	// Write to CSV
	fmt.Printf("Writing records to %s...\n", outputFile)
	if err := writeCSV(records, outputFile); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("CSV file generation completed successfully!")
}
