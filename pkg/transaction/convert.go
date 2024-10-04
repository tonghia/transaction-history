package transaction

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// readCSV reads and parses the CSV file into a slice of Transactions.
func readCSV(file *os.File) []Transaction {

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
