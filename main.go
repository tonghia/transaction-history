package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Transaction represents a single deposit or withdrawal.
type Transaction struct {
	Date    string `json:"date"`
	Amount  int    `json:"amount"`
	Content string `json:"content"`
}

// Summary represents the JSON output structure.
type Summary struct {
	Period           string        `json:"period"`
	TotalIncome      int           `json:"total_income"`
	TotalExpenditure int           `json:"total_expenditure"`
	Transactions     []Transaction `json:"transactions"`
}

func main() {
	// Set up logging to include the timestamp.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Parse command-line arguments.
	yearMonth, filePath := parseArguments()

	// Validate and parse the yearMonth input.
	year, month := parseYearMonth(yearMonth)

	// Read and parse the CSV file.
	transactions := readCSV(filePath)

	// Filter transactions based on the specified year and month.
	filteredTransactions := filterTransactions(transactions, year, month)

	// Calculate total income and expenditure.
	totalIncome, totalExpenditure := calculateTotals(filteredTransactions)

	// Sort transactions in descending order by date.
	sortTransactions(filteredTransactions)

	// Create the summary.
	summary := Summary{
		Period:           fmt.Sprintf("%04d/%02d", year, month),
		TotalIncome:      totalIncome,
		TotalExpenditure: totalExpenditure,
		Transactions:     filteredTransactions,
	}

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

// parseYearMonth validates and parses the YYYYMM input.
func parseYearMonth(yearMonth string) (int, time.Month) {
	if len(yearMonth) != 6 {
		log.Fatalf("Invalid period format. Expected YYYYMM, got: %s", yearMonth)
	}

	year, err := strconv.Atoi(yearMonth[:4])
	if err != nil {
		log.Fatalf("Invalid year in period: %v", err)
	}

	monthInt, err := strconv.Atoi(yearMonth[4:6])
	if err != nil {
		log.Fatalf("Invalid month in period: %v", err)
	}

	if monthInt < 1 || monthInt > 12 {
		log.Fatalf("Month must be between 01 and 12. Got: %02d", monthInt)
	}

	return year, time.Month(monthInt)
}

// readCSV reads and parses the CSV file into a slice of Transactions.
func readCSV(filePath string) []Transaction {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read the header line.
	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Error reading CSV header: %v", err)
	}

	// Validate headers.
	expectedHeaders := []string{"date", "amount", "content"}
	for i, header := range headers {
		if strings.TrimSpace(strings.ToLower(header)) != expectedHeaders[i] {
			log.Fatalf("Unexpected header: expected '%s', got '%s'", expectedHeaders[i], header)
		}
	}

	var transactions []Transaction

	// Read each record.
	recordNumber := 1 // Including header
	for {
		record, err := reader.Read()
		recordNumber++
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error reading CSV record at line %d: %v", recordNumber, err)
		}

		// Validate that no columns are empty.
		for i, field := range record {
			if strings.TrimSpace(field) == "" {
				log.Fatalf("Empty field in column '%s' at line %d", expectedHeaders[i], recordNumber)
			}
		}

		// Parse the date to ensure correct format.
		dateStr := record[0]
		_, err = time.Parse("2006/01/02", dateStr)
		if err != nil {
			log.Fatalf("Invalid date format at line %d: %v", recordNumber, err)
		}

		// Parse the amount.
		amount, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatalf("Invalid amount at line %d: %v", recordNumber, err)
		}

		content := record[2]

		transaction := Transaction{
			Date:    dateStr,
			Amount:  amount,
			Content: content,
		}

		transactions = append(transactions, transaction)
	}

	return transactions
}

// filterTransactions filters transactions based on the specified year and month.
func filterTransactions(transactions []Transaction, year int, month time.Month) []Transaction {
	var filtered []Transaction

	for _, tx := range transactions {
		txDate, err := time.Parse("2006/01/02", tx.Date)
		if err != nil {
			log.Fatalf("Error parsing transaction date '%s': %v", tx.Date, err)
		}

		if txDate.Year() == year && txDate.Month() == month {
			filtered = append(filtered, tx)
		}
	}

	return filtered
}

// calculateTotals calculates the total income and total expenditure.
func calculateTotals(transactions []Transaction) (int, int) {
	totalIncome := 0
	totalExpenditure := 0

	for _, tx := range transactions {
		if tx.Amount > 0 {
			totalIncome += tx.Amount
		} else {
			totalExpenditure += tx.Amount
		}
	}

	return totalIncome, totalExpenditure
}

// sortTransactions sorts transactions in descending order by date.
func sortTransactions(transactions []Transaction) {
	sort.Slice(transactions, func(i, j int) bool {
		dateI, _ := time.Parse("2006/01/02", transactions[i].Date)
		dateJ, _ := time.Parse("2006/01/02", transactions[j].Date)
		return dateI.After(dateJ)
	})
}

// outputJSON marshals the summary into JSON and outputs it.
func outputJSON(summary Summary) {
	jsonData, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	fmt.Println(string(jsonData))
}
